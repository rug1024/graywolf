<script>
  // "New message" modal launched from the sidebar/list "New" button
  // or the `?compose=1` deep link. Collects a To: callsign (via
  // CallsignAutocomplete) and delegates the actual send to the
  // parent — on success it closes and the parent navigates to the
  // newly-opened thread.
  //
  // Tactical compose skips this modal: `?compose=tactical:NET`
  // navigates straight into the tactical thread's compose bar with
  // the locked-pill To.

  import { Icon, Modal, Select } from '@chrissnell/chonky-ui';
  import CallsignAutocomplete from './CallsignAutocomplete.svelte';
  import ComposeBar from './ComposeBar.svelte';
  import RFPathPicker from './RFPathPicker.svelte';
  import { channelsStore, start as startChannels } from '../../lib/stores/channels.svelte.js';

  /** @type {{
   *    open: boolean,
   *    onSend?: (text: string, to: string, channel: number) => Promise<any>,
   *    onClose?: () => void,
   *    onPrepare?: (to: string, options: any) => Promise<any>,
   *  }}
   */
  let {
    open = $bindable(false),
    onSend,
    onPrepare,
    onClose,
  } = $props();

  let to = $state('');
  // Per-send TX channel override (issue #472). 0 = use the operator's
  // configured default channel. Tucked behind the Advanced disclosure,
  // collapsed by default so the normal compose flow stays uncluttered.
  let channel = $state(0);
  let advancedOpen = $state(false);
  let sendPath = $state('__keep');
  let rfPath = $state(null);
  let pathChanged = $state(false);
  let pathValid = $state(true);
  const sendPathOptions = [
    { value: '__keep', label: 'Keep conversation / default' },
    { value: '', label: 'Use global default' },
    { value: 'rf_only', label: 'RF only' },
    { value: 'is_only', label: 'APRS-IS only' },
    { value: 'both', label: 'RF + APRS-IS' },
  ];

  // Keep options stable for the complete batch; reset for a new session.
  let wasOpen = false;
  $effect(() => {
    if (open && !wasOpen) {
      channel = 0;
      sendPath = '__keep';
      rfPath = null;
      pathChanged = false;
      advancedOpen = false;
    }
    wasOpen = open;
  });

  // Populate the channel list only when Advanced is expanded — most
  // sends never touch it, so we don't poll on every compose.
  $effect(() => {
    if (advancedOpen) startChannels();
  });

  // Same selection semantics as Settings → Messages → Transmit channel:
  // offer APRS-eligible (non-packet) channels plus a "default" sentinel.
  const channelOptions = $derived([
    { value: 0, label: 'Default (configured TX channel)' },
    ...channelsStore.list
      .filter((c) => c.mode !== 'packet')
      .map((c) => ({ value: c.id, label: c.name })),
  ]);

  function close() {
    open = false;
    // Called once after the complete batch succeeds.
    onClose?.();
  }

  async function prepare(target) {
    // Snapshot options before awaiting any request. All fragments then use
    // the same persisted conversation policy, including the first one.
    await onPrepare?.(target, {
      sendPath: sendPath === '__keep' ? null : sendPath,
      rfPath, pathChanged,
    });
  }
  async function send(text, target, batchChannel) {
    await onSend?.(text, target.trim().toUpperCase(), batchChannel);
  }

</script>

<Modal bind:open onClose={close}>
  <Modal.Header>
    <h3 class="title">New message</h3>
    <Modal.Close aria-label="Close">
      <Icon name="x" size="lg" />
    </Modal.Close>
  </Modal.Header>
  <Modal.Body>
    <div class="to-field">
      <label for="compose-new-to">To</label>
      <CallsignAutocomplete
        bind:value={to}
        placeholder="Callsign or APRS service"
        onCommit={(v) => to = v}
        autofocus={true}
      />
    </div>
    <!-- ComposeBar in "thread" mode so it doesn't render its own
         To: field (the modal wrapper owns that above). `embedded`
         flips the compose bar from position:absolute → relative so
         the modal body handles layout. -->
    <ComposeBar
      mode="thread"
      dmPeer={to}
      threadHasMessages={false}
      autoFocus={false}
      embedded={true}
      {channel}
      sendDisabled={!pathValid}
      beforeSend={prepare}
      onSend={send}
      onSent={close}
    />
    <!-- Advanced (collapsed by default). Per-station channel selection
         for operators running multiple radios on one instance. -->
    <details class="advanced" bind:open={advancedOpen}>
      <summary class="advanced-trigger">
        <span class="chev" aria-hidden="true">{advancedOpen ? '▾' : '▸'}</span>
        Advanced
      </summary>
      <div class="advanced-body">
        <span class="advanced-label">Transmit channel</span>
        <Select
          value={channel}
          onValueChange={(v) => (channel = Number(v))}
          options={channelOptions}
          aria-label="Message transmit channel"
        />
        <p class="advanced-hint">
          Radio / interface used for RF transmission.
        </p>
        <span class="advanced-label">Send path (transport)</span>
        <Select value={sendPath} onValueChange={v => sendPath = String(v)} options={sendPathOptions} aria-label="Message send transport" />
        <p class="advanced-hint">How to send: RF = radio; APRS-IS = Internet. Automatic tries RF first, then Internet if RF is unavailable.</p>
        <RFPathPicker value={rfPath} bind:valid={pathValid} onChange={v => { rfPath = v; pathChanged = true; }} />
        <p class="advanced-hint">“RF only” selects radio transport. “RFONLY” in the packet path asks receiving iGates not to upload it. Choices apply to this conversation.</p>
      </div>
    </details>
  </Modal.Body>
</Modal>

<style>
  .title {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    font-family: var(--font-mono);
  }
  .to-field {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 16px;
  }
  .to-field label {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: var(--color-text-dim);
  }

  .advanced {
    margin-top: 12px;
    border: 1px solid var(--color-border, #ddd);
    border-radius: 6px;
  }
  .advanced-trigger {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    cursor: pointer;
    font-size: 12px;
    color: var(--color-text, #222);
    user-select: none;
    list-style: none;
  }
  /* Hide the default <details> disclosure marker on every browser. */
  .advanced-trigger::-webkit-details-marker { display: none; }
  .advanced-trigger::marker { content: ''; }
  .advanced-trigger:hover { background: var(--color-surface, #f4f4f4); }
  .advanced[open] .advanced-trigger { border-bottom: 1px solid var(--color-border, #ddd); }
  .chev {
    display: inline-block;
    width: 12px;
    color: var(--color-text-muted, #888);
    font-size: 11px;
  }
  .advanced-body {
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .advanced-label {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: var(--color-text-dim);
  }
  .advanced-hint {
    margin: 0;
    font-size: 11px;
    line-height: 1.4;
    color: var(--color-text-dim);
  }
</style>
