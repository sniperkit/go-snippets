Configuration Manager API:
Configuration Manager(CM) provides a system-independent interface for the following classes of operations:
    1. Driver Instantiation:
        a. Driver Probing: determing which driver can be used to drive a particular device.
        b. Parameter passing: passing additional run-time parameters to the driver.
    2. Hardware access:
        a. Reading/Writing hardware registers; for example, PCIe memory or I/O mapped devies registers.
        b. Registering interrup handlers.
        c. DMA-able memory allocation.
        d. Cache management.
        e. Logical/physical address translation.


The Configuaration Manager provides the abstraction for low-level functions that acess hardware to implement a SOC driver. It is a Hardware Abstraction Layer(HAL) as opposed to the System Abstraction Layer (SAL) used to abstract services provided by a typical operating system.

These abstractions are implemented in a form of function vectors associated with each of the driver instances. This makes the SOC driver highly portable and flexible with respect to various hardware configuations.

Configuration mamnager provides only the minimal number of abstractions needed to implement a SOC driver. For example, configuration mamnager does not provide any abstractions for scanning PCIe bus. This functionality is not needed to access device registers and there, in fact might not be a PCIe bus.

The soc_cm_device_vectors_t and soc_cm_dev_t strutures represent the most important abstractons provided by configuration manager.
