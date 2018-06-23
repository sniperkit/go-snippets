1.使用soc_reg_iterate(int unit, soc_reg_iter_f do_it, void *data)可以遍历chip上所有的寄存器。
2.使用parse_memory_name可以Look up a table memory by user-friendly name.
3.使用soc_mem_read可以读取一张表.
4.使用memp = &SOC_MEM_INFO (unit, mem)获取表的内容, 然后就可以遍历该表的所有field了。
5.使用soc_mem_field_get(unit, mem, entry, field, fval)来获取某张表的某一域的值。
6.soc_driver_t 定义了chip的驱动，其中包括寄存器定义，表定义，port, feature, counter以及chip特定初始化历程等。
7. SOC_PORT_NAME(unit, (port))可以根据端口号获取端口的名字
