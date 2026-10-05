package com.nw5w.graywolf.platformsvc

import android.content.Context
import com.google.protobuf.ByteString
import com.nw5w.graywolf.platformproto.PlatformMessage
import com.nw5w.graywolf.platformproto.SerialClose
import com.nw5w.graywolf.platformproto.SerialData
import com.nw5w.graywolf.platformproto.SerialKind
import com.nw5w.graywolf.platformproto.SerialOpen
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertArrayEquals
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test
import org.mockito.Mockito.mock

class BleAdapterTest {
    @Test fun scanRequest_forwardsDiscoveredDevice() = runTest {
        val dispatcher = StandardTestDispatcher(testScheduler)
        val facade = FakeBleFacade(
            scanDevices = listOf(BleFoundDevice("AA:BB:CC:DD:EE:01", "T-Beam", -47))
        )
        val sent = mutableListOf<PlatformMessage>()
        val adapter = BleAdapter(facade, mock(Context::class.java), dispatcher) { sent.add(it) }

        adapter.handleScanRequest()

        val results = sent.filter { it.hasBleScanResult() }.map { it.bleScanResult }
        assertEquals(1, results.size)
        assertEquals("AA:BB:CC:DD:EE:01", results[0].addr)
        assertEquals("T-Beam", results[0].name)
        assertEquals(-47, results[0].rssi)
        assertTrue(facade.scanning)

        adapter.handleScanStop()
        assertFalse(facade.scanning)
    }

    @Test fun serialOpen_wrongKind_repliesWithError() = runTest {
        val dispatcher = StandardTestDispatcher(testScheduler)
        val facade = FakeBleFacade()
        val sent = mutableListOf<PlatformMessage>()
        val adapter = BleAdapter(facade, mock(Context::class.java), dispatcher) { sent.add(it) }

        adapter.handleSerialOpen(
            SerialOpen.newBuilder()
                .setHandle(7)
                .setKind(SerialKind.SERIAL_KIND_BLUETOOTH)
                .setAddress("AA:BB:CC:DD:EE:02")
                .build()
        )

        val ack = sent.single { it.hasSerialOpenAck() }.serialOpenAck
        assertEquals(7, ack.handle)
        assertFalse(ack.ok)
        assertTrue(ack.error.contains("unsupported_kind"))
    }

    @Test fun bleSerialOpen_write_andDisconnect_preserveStreamLifecycle() = runTest {
        val dispatcher = StandardTestDispatcher(testScheduler)
        val session = FakeBleGattSession()
        val facade = FakeBleFacade(session = session)
        val sent = mutableListOf<PlatformMessage>()
        val adapter = BleAdapter(facade, mock(Context::class.java), dispatcher) { sent.add(it) }

        adapter.handleSerialOpen(
            SerialOpen.newBuilder()
                .setHandle(42)
                .setKind(SerialKind.SERIAL_KIND_BLE)
                .setAddress("AA:BB:CC:DD:EE:03")
                .build()
        )
        advanceUntilIdle()

        val ack = sent.single { it.hasSerialOpenAck() }.serialOpenAck
        assertTrue(ack.ok)
        assertEquals("AA:BB:CC:DD:EE:03", facade.openedMac)

        val payload = byteArrayOf(0xC0.toByte(), 0x00, 0x01, 0xC0.toByte())
        adapter.handleSerialData(
            SerialData.newBuilder()
                .setHandle(42)
                .setData(ByteString.copyFrom(payload))
                .build()
        )
        advanceUntilIdle()
        assertEquals(1, session.writes.size)
        assertArrayEquals(payload, session.writes[0])

        val inbound = byteArrayOf(0xC0.toByte(), 0x00, 0x55, 0xC0.toByte())
        session.emit(inbound)
        advanceUntilIdle()
        val rx = sent.single { it.hasSerialData() }.serialData.data.toByteArray()
        assertArrayEquals(inbound, rx)

        session.disconnect()
        advanceUntilIdle()
        val closes = sent.filter { it.hasSerialClose() }.map { it.serialClose }
        assertEquals(1, closes.size)
        assertEquals(42, closes[0].handle)
        assertTrue(closes[0].reason.contains("ble_link_lost"))
        assertTrue(session.closed)
    }

    @Test fun closeUnknownHandle_isNoOp() = runTest {
        val dispatcher = StandardTestDispatcher(testScheduler)
        val sent = mutableListOf<PlatformMessage>()
        val adapter = BleAdapter(
            FakeBleFacade(),
            mock(Context::class.java),
            dispatcher,
        ) { sent.add(it) }

        adapter.handleSerialClose(
            SerialClose.newBuilder().setHandle(999).setReason("test").build()
        )
        advanceUntilIdle()
        assertTrue(sent.isEmpty())
    }

    private class FakeBleFacade(
        private val scanDevices: List<BleFoundDevice> = emptyList(),
        private val session: FakeBleGattSession = FakeBleGattSession(),
    ) : BleFacade {
        var scanning = false
        var openedMac: String? = null

        override fun startScan(
            onFound: (BleFoundDevice) -> Unit,
            onError: ((String) -> Unit)?,
        ) {
            scanning = true
            scanDevices.forEach(onFound)
        }

        override fun stopScan() {
            scanning = false
        }

        override fun openGatt(context: Context, mac: String): BleGattSession {
            openedMac = mac
            return session
        }

        override fun removeBond(mac: String): Boolean = true
    }

    private class FakeBleGattSession : BleGattSession {
        val writes = mutableListOf<ByteArray>()
        var closed = false
        private var dataCallback: ((ByteArray) -> Unit)? = null
        private var disconnectCallback: (() -> Unit)? = null

        override fun write(bytes: ByteArray) {
            writes += bytes.copyOf()
        }

        override fun onData(cb: (ByteArray) -> Unit) {
            dataCallback = cb
        }

        override fun onDisconnect(cb: () -> Unit) {
            disconnectCallback = cb
        }

        override fun close() {
            closed = true
        }

        fun emit(bytes: ByteArray) {
            dataCallback?.invoke(bytes.copyOf())
        }

        fun disconnect() {
            disconnectCallback?.invoke()
        }
    }
}
