package com.nw5w.graywolf.audio

import android.content.Context
import android.database.sqlite.SQLiteDatabase
import android.util.Log
import java.io.File

/**
 * Decides whether the current persisted configuration actually needs Android
 * microphone capture. KISS-only channels have a NULL input_device_id and must
 * not cause RECORD_AUDIO to be requested or AudioRecord to be opened.
 *
 * The legacy audio_device_id probe keeps upgrades safe before the Go-side
 * migration has had a chance to rename the column. Any unreadable/missing DB
 * is treated as no configured audio input; a fresh install is KISS-capable
 * without granting microphone access.
 */
object AudioConfigGate {
    private const val TAG = "AudioConfigGate"

    fun requiresMicrophone(context: Context): Boolean {
        val dbFile = File(context.filesDir, "graywolf.db")
        if (!dbFile.exists()) return false

        var db: SQLiteDatabase? = null
        return try {
            db = SQLiteDatabase.openDatabase(
                dbFile.absolutePath,
                null,
                SQLiteDatabase.OPEN_READONLY,
            )
            val columns = mutableSetOf<String>()
            db.rawQuery("PRAGMA table_info(channels)", null).use { cursor ->
                val nameIndex = cursor.getColumnIndex("name")
                while (cursor.moveToNext() && nameIndex >= 0) {
                    columns += cursor.getString(nameIndex)
                }
            }
            val inputColumn = when {
                "input_device_id" in columns -> "input_device_id"
                "audio_device_id" in columns -> "audio_device_id"
                else -> return false
            }
            val enabledClause = if ("enabled" in columns) "enabled = 1 AND " else ""
            db.rawQuery(
                "SELECT 1 FROM channels WHERE ${enabledClause}${inputColumn} IS NOT NULL LIMIT 1",
                null,
            ).use { it.moveToFirst() }
        } catch (t: Throwable) {
            Log.w(TAG, "could not inspect audio configuration; not starting microphone capture", t)
            false
        } finally {
            try { db?.close() } catch (_: Throwable) {}
        }
    }
}
