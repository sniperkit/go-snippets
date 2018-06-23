1. Basic:
    bcm debug bcmapi l3 +
    bcm debug bcmapi +
    bcm memwatch vlan read
    bcm memwatch vlan write
    bcm memwatch vlan off
    bcm lls
    bcm cos bandwith_show
    bcm cos discard_show
    bcm cos port show PBM=all
    bcm d sw cosq
    bcm cos show
    bcm show pmap
    bcm dump sw port
    bcm dump sw vlan
    bcm dump sw ifp
    bcm fp show 
    bcm show features
    bcm show st
    bcm show params
    bcm hash
    bcm HashDestination
    bcm pr
    bcm pp
    bcm rate
    bcm knetctrl netif show 
    bcm knetctrl filter show
    bcm egress show
2.  TechSupport
    The "techsupport" utility internally maintains the following sets of data on a per chip per feature basis:

        diag shell commands (diag).
        Register list (reg).
        Memory table list (mem).
        Software maintained feature specific state (sw), where applicable. 
    The user can collect the data described above using the techsupport command. This command has the intelligence to identify the underlying chip on which it is being run and collects the chip specific data.

    The following is the format of techsupport command:

        techsupport basic - collects basic config/setup information from switch
        techsupport [diag] [reg] [mem] [list] [sw] [verbose]
    Command "techsupport basic" collects basic config/setup information from switch. It executes the following commands and is supported on all XGS chips.

        attach
        version
        config show
        show counters
        linkscan
        ps
        lls
        hsp
        soc
        fp show
        show pmap
        phy info
        show params
        show features
        dump soc diff
        dump pcic
   
    The command "techsupport " dumps all four sets of feature specific information. By default (i.e if verbose option is not mentioned), memory and register dumps are changes from defaults. If any other options (diag, reg, mem, sw) are specified, then only that specific information is dumped. 

    For example,
        Command "techsupport l3mc diag" executes the Layer 3 multicast related diag shell commands.
        Command "techsupport l3mc diag mem" executes the Layer3 multicast related diag shell commands and dumps the Layer3 multicast related memory tables. 
        Command "techsupport l3mc"  collects all the above 4 sets of data.
        Command "techsupport l3mc list" lists the diag shell commands, memory and register names.
    As of SDK 6.5.7, the techsupport command supports the following features. These features are supported for Trident2, Trident2+ and Tomahawk chipsets.
        l3uc - collects L3 unicast related debug information
        l3mc - collects L3 multicast related debug information
        mpls - collects mpls related debug information
        mmu - collects mmu related debug information
        niv - collects niv related debug information
        riot - collects riot related debug information
        oam - collects oam related debug information
        vxlan - collects vxlan related debug information
        vlan - collects vlan related debug information
        cos - collects cos related debug information
        load-balance - collects RTAG7,DLB,ECMP and trunk related debug information   
        efp - collects efp related debug information
        ifp - collects ifp related debug information
        vfp - collects vfp related debug information

