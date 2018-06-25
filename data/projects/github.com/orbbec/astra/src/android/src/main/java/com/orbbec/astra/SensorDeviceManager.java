package com.orbbec.astra;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.hardware.usb.UsbDevice;
import android.hardware.usb.UsbManager;
import android.util.Log;

import java.util.*;

public class SensorDeviceManager {

    private static final String LOG_TAG = "SensorDeviceManager";
    private final Context androidContext;
    private final UsbManager usbManager;
    private final UsbDeviceAccessBroker usbAccessBroker;
    private final HashSet<UsbDeviceInfo> knownDeviceTypes;
    private final HashSet<UsbDevice> availableDevices = new HashSet<UsbDevice>();

    private final Queue<UsbDevice> requestQueue = new LinkedList<UsbDevice>();

    private boolean isOpeningDevices = false;
    private SensorDeviceManagerListener listener;

    private final UsbDeviceAccessBroker.UsbDeviceAccessEventListener accessListener =
            new UsbDeviceAccessBroker.UsbDeviceAccessEventListener() {

                @Override
                public void onDeviceOpened(UsbDevice device) {
                    availableDevices.add(device);
                    onDeviceOpenCallback();
                }

                @Override
                public void onDeviceOpenFailed(UsbDevice device) {
                    onDeviceOpenCallback();
                }

                private void onDeviceOpenCallback()
                {
                    if (requestQueue.isEmpty()) {
                        isOpeningDevices = false;
                        if (listener != null)
                            listener.onOpenAllDevicesCompleted(availableDevices);
                    } else {
                        usbAccessBroker.open(requestQueue.remove(), accessListener);
                    }
                }
            };



    public Iterable<UsbDevice> getAvailableDevices() { return availableDevices; }

    public SensorDeviceManager(Context androidContext, SensorDeviceManagerListener listener) {
        this.androidContext = androidContext;
        this.usbManager = (UsbManager)androidContext.getSystemService(Context.USB_SERVICE);
        this.usbAccessBroker = new UsbDeviceAccessBroker(androidContext);
        this.listener = listener;

        knownDeviceTypes = getKnownDeviceTypes();

        IntentFilter filter = new IntentFilter(AstraDeviceMonitorActivity.ACTION_USB_DEVICE_ATTACHED);
        androidContext.registerReceiver(new BroadcastReceiver() {
            @Override
            public void onReceive(Context context, Intent intent) {
                if (!intent.getAction().equals(AstraDeviceMonitorActivity.ACTION_USB_DEVICE_ATTACHED))
                    return;

                UsbDevice usbDevice = (UsbDevice)intent.getParcelableExtra(UsbManager.EXTRA_DEVICE);

                if (!availableDevices.contains(usbDevice)) {
                    openDevice(usbDevice);
                }

            }}, filter);
    }

    public void openDevice(UsbDevice device) {

        if (device == null)
            throw new IllegalArgumentException("device");

        usbAccessBroker.open(device, new UsbDeviceAccessBroker.UsbDeviceAccessEventListener() {
            @Override
            public void onDeviceOpened(UsbDevice device) {
                availableDevices.add(device);
                if (listener != null) {
                    listener.onOpenDeviceCompleted(device, true);
                }
            }

            @Override
            public void onDeviceOpenFailed(UsbDevice device) {
                if (listener != null) {
                    listener.onOpenDeviceCompleted(device, false);
                }
            }
        });
    }

    public void openAllDevices() {

        if (isOpeningDevices)
            return;

        Log.d(LOG_TAG, "opening all devices");
        isOpeningDevices = true;

        for(UsbDevice device : getConnectedSensorDevices()) {
            requestQueue.add(device);
        }

        if (requestQueue.isEmpty())
        {
            Log.d(LOG_TAG, "no devices to open");
        } else {
            usbAccessBroker.open(requestQueue.remove(), accessListener);
        }
    }

    private static HashSet<UsbDeviceInfo> getKnownDeviceTypes() {

        HashSet<UsbDeviceInfo> knownDeviceTypes = new HashSet<UsbDeviceInfo>();

        //TODO: load these from a config file as well
        //ASUS Xtion Pro
        knownDeviceTypes.add(new UsbDeviceInfo(0x1d27, 0x0600));

        //ASUS Xtion Pro Live
        knownDeviceTypes.add(new UsbDeviceInfo(0x1d27, 0x0601));

        //Orbbec Astra
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0400));
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0401));
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0402));
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0403));
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0404));
        knownDeviceTypes.add(new UsbDeviceInfo(0x2bc5, 0x0405));

        return knownDeviceTypes;
    }

    private Iterable<UsbDevice> getConnectedSensorDevices()
    {
        ArrayList<UsbDevice> availableDevices = new ArrayList<UsbDevice>();

        for(UsbDevice device : usbManager.getDeviceList().values()) {

            UsbDeviceInfo info = new UsbDeviceInfo(device.getVendorId(), device.getProductId());
            if (knownDeviceTypes.contains(info))
                availableDevices.add(device);
        }

        return availableDevices;
    }

}
