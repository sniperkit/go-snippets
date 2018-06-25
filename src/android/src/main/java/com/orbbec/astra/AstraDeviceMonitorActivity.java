package com.orbbec.astra;

import android.app.Activity;
import android.content.Intent;
import android.hardware.usb.UsbManager;
import android.os.Bundle;
import android.os.Parcelable;

public class AstraDeviceMonitorActivity extends Activity {

    public static final String ACTION_USB_DEVICE_ATTACHED = "com.orbbec.ACTION_USB_DEVICE_ATTACHED";

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
    }

    @Override
    protected void onResume() {
        super.onResume();

        Intent intent = getIntent();
        if (intent == null) return;

        String action = intent.getAction();
        if (!UsbManager.ACTION_USB_DEVICE_ATTACHED.equals(action))
            return;

        Parcelable usbDeviceParcel = intent.getParcelableExtra(UsbManager.EXTRA_DEVICE);

        Intent broadcastIntent = new Intent(ACTION_USB_DEVICE_ATTACHED);
        broadcastIntent.putExtra(UsbManager.EXTRA_DEVICE, usbDeviceParcel);

        sendBroadcast(broadcastIntent);
    }
}