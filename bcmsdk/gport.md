1. Traditionally, APIS have had multiple parameters to specify destination information(such as module, port, or trunk).
    The GPORT type was introduced to provide a single 32-bit identifier for destinations. The following GPORT destination
    types are provided:
        1. Local port number (a physical port on the local unit)
        2. {moulde, port} pair
        3. Trunk identifier
        4. Black Hole
        5. Local CPU.
        6. L2 MPLS virtual-port (VPLS/VPWS)
        7. Subport group.
        8. Subport virtual-port
        9. MAC-in-MAC virtual-port
        10. WLAN virutal-port
        11. CoSQ scheduler
        12. Tunnel identifier
        13. {device ID, port} pair

