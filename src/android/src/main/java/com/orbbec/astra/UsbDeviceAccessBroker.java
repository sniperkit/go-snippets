package com.orbbec.astra;

import android.app.PendingIntent;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.hardware.usb.UsbDevice;
import android.hardware.usb.UsbManager;
import android.util.Log;

import java.util.Hashtable;

public class UsbDeviceAccessBroker {

    private static final String LOG_TAG = "UsbDeviceAccessBroker";

    public interface UsbDeviceAccessEventListener {
        void onDeviceOpened(UsbDevice device);
        void onDeviceOpenFailed(UsbDevice device);
    }

    private static final String LISTENER_ID = "com.orbbec.LISTENER_ID";

    private final String USB_PERMISSION_INTENT;
    private final UsbManager usbManager;
    private final Context androidContext;

    private final Hashtable<Integer, UsbDeviceAccessEventListener> listeners =
            new Hashtable<Integer, UsbDeviceAccessEventListener>();

    private static int intentId = 0;
    private static synchronized int getNextId() { return intentId++; }

    private final BroadcastReceiver permissionIntentReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {

            String action = intent.getAction();
            if (!USB_PERMISSION_INTENT.equals(action))
                return;

            synchronized (this) {

                int intentInstanceId = intent.getIntExtra(LISTENER_ID, -1);

                UsbDeviceAccessEventListener listener = null;

                if (intentInstanceId != -1 && listeners.containsKey(intentInstanceId)) {
                    listener = listeners.get(intentInstanceId);
                    listeners.remove(intentInstanceId);
                }

                UsbDevice device = intent.getParcelableExtra(UsbManager.EXTRA_DEVICE);

                boolean permissionGranted = intent.getBooleanExtra(UsbManager.EXTRA_PERMISSION_GRANTED, false);

                Log.d(LOG_TAG,
                        "received permission request result for "
                                + device.getDeviceName()
                                + ", granted: " + permissionGranted);

                if (device == null && listener != null) {
                    Log.d(LOG_TAG, "notifying listener of device open failure");
                    listener.onDeviceOpenFailed(null);
                }

                if (permissionGranted) {
                    Log.d(LOG_TAG, "opening device: " + device.getDeviceName());
                    usbManager.openDevice(device);
                    Log.d(LOG_TAG, "opened device: " + device.getDeviceName());

                    if (listener != null) {
                        Log.d(LOG_TAG, "notifying listener of device open success");
                        listener.onDeviceOpened(device);
                    }
                } else {
                    if (listener != null) {
                        Log.d(LOG_TAG, "notifying listener of device open failure");
                        listener.onDeviceOpenFailed(device);
                    }

                }
            }
        }
    };

    public UsbDeviceAccessBroker(Context androidContext) {
        this.androidContext = androidContext;
        this.usbManager = (UsbManager)androidContext.getSystemService(Context.USB_SERVICE);
        this.USB_PERMISSION_INTENT = androidContext.getPackageName() + ".USB_PERMISSION";

        IntentFilter filter = new IntentFilter(USB_PERMISSION_INTENT);
        androidContext.registerReceiver(permissionIntentReceiver, filter);
    }

    public void open(UsbDevice device, UsbDeviceAccessEventListener listener) {

        if (device == null)
            throw new IllegalArgumentException("device");

        Log.d(LOG_TAG,"attempting to open: " + device.getDeviceName());

        if (usbManager.hasPermission(device)) {
            Log.d(LOG_TAG, "opening device: " + device.getDeviceName());
            usbManager.openDevice(device);
            Log.d(LOG_TAG, "opened device: " + device.getDeviceName());

            Log.d(LOG_TAG, "notifying listener of device open success");
            listener.onDeviceOpened(device);
            return;
        }

        Intent permissionIntent = new Intent(USB_PERMISSION_INTENT);

        if (listener != null) {
            int id = getNextId();
            Log.d(LOG_TAG, "storing listener ("+ id + ")");
            permissionIntent.putExtra(LISTENER_ID, id);

            listeners.put(id, listener);
        }

        Log.d(LOG_TAG, "executing permission request for device: " + device.getDeviceName());

        PendingIntent intent = PendingIntent.getBroadcast(androidContext, 0, permissionIntent, 0);
        usbManager.requestPermission(device, intent);
    }
}
