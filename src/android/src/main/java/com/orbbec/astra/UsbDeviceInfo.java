package com.orbbec.astra;

public class UsbDeviceInfo {

    private final int vendorId;
    private final int productId;

    public int getVendorId() {
        return vendorId;
    }

    public int getProductId() {
        return productId;
    }

    public UsbDeviceInfo(int vendorId, int productId) {
        this.vendorId = vendorId;
        this.productId = productId;
    }

    @Override
    public String toString() {
        return "UsbDeviceInfo{" +
                "vendorId=" + vendorId +
                ", productId=" + productId +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        UsbDeviceInfo that = (UsbDeviceInfo) o;

        if (vendorId != that.vendorId) return false;
        return productId == that.productId;

    }

    @Override
    public int hashCode() {
        int result = vendorId;
        result = 31 * result + productId;
        return result;
    }
}
