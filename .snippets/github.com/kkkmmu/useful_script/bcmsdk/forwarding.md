1. L3_IIF: Layer 3 Input Interface Properties
    每一个L3 interface创建以后都会在该表中创建一个表项。该表主要用于对三层接口入方向的控制。比如说三层接口的VRF就配置在该表中。
    同时该表项中会有一个到L3_IIF_PROFILE的索引，而L3_IIF_PROFILE中包含了对三层接口更广泛的控制功能。比如说URPF，三层转发使能，三层组播使能，ICMP REDIRECT上送CPU，以及该接口是否允许GLOBAL routing.
    在TD3设备上碰到问题：所有L3 转发表项配置均正确，但是三层流量还是上送CPU，最后找到原因是该三层接口对应的L3_IIF_PROFILE中，不允许GLOBAL routing，而我们默认配置的路由表项均为Global.
