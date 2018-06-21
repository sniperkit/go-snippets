#include <linux/module.h>
#include <linux/version.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kprobes.h>
#include <linux/kallsyms.h>

static unsigned int counter = 0;

static long pre_do_fork(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_do_fork: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}

static long pre_schedule(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_schedule: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}

static long pre_ip_rcv(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_ip_rcv: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}
static long pre_netif_rx(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_netif_rx: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}

static long pre_dev_queue_xmit(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_dev_queue_xmit: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}

static long pre_switch_to(unsigned long clone_flags, unsigned long stack_start,
			struct pt_regs *regs, unsigned long stack_size, int __user *parent_tidptr,
				int __user *child_tidptr)
{
	    printk(KERN_EMERG "pre_switch_to: counter: %u\n", counter++);
	        jprobe_return();
		    return 0;
}

static struct jprobe jprobe_do_fork = 
{
	    .entry = pre_do_fork,
	        .kp = {
				.symbol_name = "do_fork",
				    },
};

static struct jprobe jprobe_netif_rx = 
{
	    .entry = pre_netif_rx,
	        .kp = {
				.symbol_name = "netif_rx",
				    },
};

static struct jprobe jprobe_schedule = 
{
	    .entry = pre_schedule,
	        .kp = {
				.symbol_name = "schedule",
				    },
};

static struct jprobe jprobe_switch_to = 
{
	    .entry = pre_switch_to,
	        .kp = {
				.symbol_name = "__switch_to",
				    },
};

static struct jprobe jprobe_ip_rcv = 
{
	    .entry = pre_ip_rcv,
	        .kp = {
				.symbol_name = "ip_rcv",
				    },
};

static struct jprobe jprobe_dev_queue_xmit = 
{
	    .entry = pre_dev_queue_xmit,
	        .kp = {
				.symbol_name = "dev_queue_xmit",
				    },
};

static int __init jp_init(void)
{
	    int ret = -1;
	        printk("Test jp module init\n");

		    ret = register_jprobe(&jprobe_do_fork);
		        ret = register_jprobe(&jprobe_netif_rx);
			    ret = register_jprobe(&jprobe_ip_rcv);
			        ret = register_jprobe(&jprobe_schedule);
				    ret = register_jprobe(&jprobe_switch_to);
				        ret = register_jprobe(&jprobe_dev_queue_xmit);
					    if (ret < 0)
						        {
									printk (KERN_EMERG "Register jprobe failed, return: %x\n", ret);
										return -1;
										    }

#if 0
					        printk(KERN_EMERG "Planted jprobe at: %p, handler addr: %p\n", 
									    jprobe_do_fork.kp.addr, jprobe_do_fork.entry);
#endif

						    return 0;
}

static void __exit jp_exit(void)
{
	    unregister_jprobe(&jprobe_do_fork);
	        unregister_jprobe(&jprobe_netif_rx);
		    unregister_jprobe(&jprobe_ip_rcv);
		        unregister_jprobe(&jprobe_dev_queue_xmit);
			    printk("Test jp module removed\n");
}

module_init(jp_init);
module_exit(jp_exit);

MODULE_AUTHOR("kkkmmu");
MODULE_DESCRIPTION("Jprobe_Mudule");
MODULE_LICENSE("GPL");

