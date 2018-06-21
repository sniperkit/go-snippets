#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kallsyms.h>
#include <linux/perf_event.h>
#include <linux/hw_breakpoint.h>

struct perf_event* __percpu* ip_rcv_hbp;
struct perf_event* __percpu* init_net_hbp;
struct perf_event* __percpu* dev_queue_xmit_hbp;

static int tttcount = 0;
//static char ksym_name[KSYM_NAME_LEN] = "pid_max";
//module_param_string(ksym, ksym_name, KSYM_NAME_LEN, S_IRUGO);

#if 0
MODULE_PARM_DES(ksym, "Kernel symbol to monitor; this module will report any ");
#endif

static void hw_ip_rcv(struct perf_event *bp, 
			struct perf_sample_data *data,
				struct pt_regs* regs)
{
	    printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
	        printk(KERN_EMERG "	    ip_rcv has been changed: %d\n", tttcount++);
		    dump_stack();
		        printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
}

static void hw_dev_queue_xmit(struct perf_event *bp, 
			struct perf_sample_data *data,
				struct pt_regs* regs)
{
	    printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
	        printk(KERN_EMERG "	    dev_queue_xmit has been changed: %d\n", tttcount++);
		    dump_stack();
		        printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
}

static void hw_init_net(struct perf_event *bp, 
			struct perf_sample_data *data,
				struct pt_regs* regs)
{
	    printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
	        printk(KERN_EMERG "	    init_net_hbp has been changed: %d\n", tttcount++);
		    dump_stack();
		        printk(KERN_EMERG "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n");
}

static int __init hw_break_module_init(void)
{
	    int ret;
	        struct perf_event_attr attr;
		    struct perf_event_attr attr1;
		        struct perf_event_attr attr2;

			    hw_breakpoint_init(&attr);
			        attr.bp_addr = kallsyms_lookup_name("ip_rcv");
				    //attr.bp_addr = kallsyms_lookup_name("__current_thread_info->task");
				    attr.bp_len = HW_BREAKPOINT_LEN_1;
				        attr.bp_type = HW_BREAKPOINT_W;

					    hw_breakpoint_init(&attr1);
					        attr.bp_addr = kallsyms_lookup_name("init_net");
						    //attr.bp_addr = kallsyms_lookup_name("__current_thread_info->task");
						    attr.bp_len = HW_BREAKPOINT_LEN_4;
						        attr.bp_type = HW_BREAKPOINT_W | HW_BREAKPOINT_R;

							    hw_breakpoint_init(&attr2);
							        attr.bp_addr = kallsyms_lookup_name("dev_queue_xmit");
								    //attr.bp_addr = kallsyms_lookup_name("__current_thread_info->task");
								    attr.bp_len = HW_BREAKPOINT_LEN_1;
								        attr.bp_type = HW_BREAKPOINT_W;

									    ip_rcv_hbp = register_wide_hw_breakpoint(&attr, hw_ip_rcv, NULL);
									        if (IS_ERR((void __force *)ip_rcv_hbp)) {
												ret = PTR_ERR((void __force *)ip_rcv_hbp);
													goto fail;
													    }

										    dev_queue_xmit_hbp = register_wide_hw_breakpoint(&attr, hw_dev_queue_xmit, NULL);
										        if (IS_ERR((void __force *)dev_queue_xmit_hbp)) {
													ret = PTR_ERR((void __force *)dev_queue_xmit_hbp);
														goto fail;
														    }

											    init_net_hbp = register_wide_hw_breakpoint(&attr, hw_init_net, NULL);
											        if (IS_ERR((void __force *)init_net_hbp)) {
														ret = PTR_ERR((void __force *)init_net_hbp);
															goto fail;
															    }

												    printk(KERN_EMERG "Hw Breakpoint for write installed!\n");

												        return 0;

fail:
													    printk(KERN_EMERG "Breakpoint registration failed!\n");

													        return ret;
}

static void __exit hw_break_module_exit(void)
{
	    unregister_wide_hw_breakpoint(ip_rcv_hbp);
	        unregister_wide_hw_breakpoint(init_net_hbp);
		    unregister_wide_hw_breakpoint(dev_queue_xmit_hbp);
		        printk(KERN_EMERG "HW breakpoint for write uninstalled!\n");
}

module_init(hw_break_module_init);
module_exit(hw_break_module_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("kkkmmu");
MODULE_DESCRIPTION("Ksym breakpoint");

