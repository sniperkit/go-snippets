package com.orbbec.astra;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.hardware.usb.UsbDevice;
import android.util.Log;

import com.orbbec.jni.Astra;

public class AstraContext {

    private static final String LOG_TAG = "AstraContext";
    private static final String RESOURCE_URI = "RESOURCE_URI";

    private final Context androidContext;
    private final SensorDeviceManager deviceManager;
    private final String RESOURCE_AVAILABLE_INTENT;


    public AstraContext(Context androidContext) {
        this.androidContext = androidContext;
        this.RESOURCE_AVAILABLE_INTENT = androidContext.getPackageName() + ".RESOURCE_AVAILABLE";
        this.deviceManager = new SensorDeviceManager(androidContext,
                new SensorDeviceManagerListener() {
                    @Override
                    public void onOpenAllDevicesCompleted(Iterable<UsbDevice> availableDevices) {
                        finishInitialization(availableDevices);
                    }

                    @Override
                    public void onOpenDeviceCompleted(UsbDevice device, boolean opened) {
                        if (opened) {
                            onNewDeviceAvailable(device);
                        }
                    }
                });

        IntentFilter filter = new IntentFilter(RESOURCE_AVAILABLE_INTENT);
        this.androidContext.registerReceiver(resourceAvailableReceiver, filter);
    }

    private void onNewDeviceAvailable(UsbDevice device) {
        Log.d(LOG_TAG, "Available device: " + device.getDeviceName());
    }

    private void finishInitialization(Iterable<UsbDevice> availableDevices) {
        for(UsbDevice device : availableDevices) {
            Log.d(LOG_TAG, "Available device: " + device.getDeviceName());

            String deviceName = device.getDeviceName();
            String[] pathComponents = deviceName.substring(1).split("/");

            if (pathComponents.length == 5
                    && pathComponents[2].equals("usb"))
            {
                int bus = Integer.parseInt(pathComponents[3], 10);
                int address = Integer.parseInt(pathComponents[4], 10);

                String resourceUri = createAstraResourceUri(
                        device.getVendorId(),
                        device.getProductId(),
                        bus,
                        address);

                Log.d(LOG_TAG, "device URI: " + resourceUri);
                Astra.notify_resource_available(resourceUri);
            }
        }
        Log.d(LOG_TAG, "Finished initialization");
    }

    private String createAstraResourceUri(int vendorId, int productId, int bus, int address)
    {
        return "usb/" + vendorId + "/" + productId + "/" + bus + "/" + address;
    }

    public void initialize() {
        Log.d(LOG_TAG, "initializing");
        Astra.initialize();
        deviceManager.openAllDevices();
    }

    private final BroadcastReceiver resourceAvailableReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {
            String action = intent.getAction();

            if (!RESOURCE_AVAILABLE_INTENT.equals(action)) return;

            String resourceUri = intent.getStringExtra(AstraContext.RESOURCE_URI);

            synchronized (this) {
                Log.d(LOG_TAG, "Resource Available: " + resourceUri);
            }
        }
    };
}
