package igate

import (
    "testing"

    "github.com/chrissnell/graywolf/pkg/aprs"
    "github.com/chrissnell/graywolf/pkg/igate/filters"
)

func TestEncodeTNC2PreservesRepeatedDigipeaterHBit(t *testing.T) {
    path := []string{"DB0KX-2*", "WIDE2-1"}
    raw := buildRawFrame(t, "DB0OTV-2", "APOT21", path, "!4903.50N/07201.75W>Test")
    pkt := &aprs.DecodedAPRSPacket{
        Source: "DB0OTV-2", Dest: "APOT21", Path: path, Raw: raw,
    }

    for _, q := range []string{"qAR", "qAO"} {
        t.Run(q, func(t *testing.T) {
            got, err := encodeTNC2(pkt, "DA6PG-1", q)
            if err != nil {
                t.Fatal(err)
            }
            want := "DB0OTV-2>APOT21,DB0KX-2*,WIDE2-1," + q + ",DA6PG-1:!4903.50N/07201.75W>Test"
            if got != want {
                t.Fatalf("TNC2 output:\n got %q\nwant %q", got, want)
            }
        })
    }
    if path[0] != "DB0KX-2*" {
        t.Fatal("original RF path was modified")
    }
}

func TestEncodeTNC2RejectsUnsupportedQConstruct(t *testing.T) {
    raw := buildRawFrame(t, "W5ABC-7", "APRS", nil, ">Test")
    pkt := &aprs.DecodedAPRSPacket{Source: "W5ABC-7", Dest: "APRS", Raw: raw}
    if _, err := encodeTNC2(pkt, "KE7XYZ", "qAC"); err == nil {
        t.Fatal("expected rejected unsupported q construct")
    }
}

func TestQConstructForRFDeliveryEligibility(t *testing.T) {
    allowMessages := []filters.Rule{
        {ID: 1, Priority: 1, Type: filters.TypeMessageDest, Pattern: "*", Action: filters.Allow},
    }
    cases := []struct {
        name string
        rxChannel int
        txChannel uint32
        path []string
        enableTx bool
        allowFilter bool
        want string
    }{
        {"direct and downlink on same channel", 7, 7, []string{"WIDE2-1"}, true, true, "qAR"},
        {"receive-only", 7, 7, nil, false, true, "qAO"},
        {"no permitted message rule", 7, 7, nil, true, false, "qAO"},
        {"received via digipeater", 7, 7, []string{"DB0KX-2*", "WIDE2-1"}, true, true, "qAO"},
        {"received on different radio channel", 8, 7, nil, true, true, "qAO"},
        {"unknown ingress channel", 0, 7, nil, true, true, "qAO"},
        {"no selected transmit channel", 7, 0, nil, true, true, "qAO"},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            cfg := Config{
                Server: "127.0.0.1:1", StationCallsign: "DA6PG-1",
                TxChannel: tc.txChannel,
            }
            if tc.enableTx { cfg.Governor = &stubGovernor{} }
            if tc.allowFilter { cfg.Rules = allowMessages }
            ig, err := New(cfg)
            if err != nil { t.Fatal(err) }
            pkt := &aprs.DecodedAPRSPacket{
                Source: "DL5YCU-14", Dest: "APRS",
                Channel: tc.rxChannel, Path: tc.path,
            }
            if got := ig.qConstructFor(pkt); got != tc.want {
                t.Fatalf("qConstructFor = %s, want %s", got, tc.want)
            }
        })
    }
}
