#include <linux/module.h>
#include <linux/version.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kprobes.h>
#include <linux/kallsyms.h>

static unsigned int counter = 0;

int pre_handler(struct kprobe* kp, struct pt_regs* regs)
{
	    printk(KERN_EMERG "pr_handler: counter: %u\n", counter++);
	        return 0;
}

void post_handler(struct kprobe* kp, struct pt_regs* regs, unsigned long flags)
{
	    printk(KERN_EMERG "post_handler: counter=%u\n", counter++);
}

static struct kprobe kp;

static int __init kp_init(void)
{
	    int ret = -1;
	        printk("Test kp module init\n");
		    kp.pre_handler = pre_handler;
		        kp.post_handler = post_handler;
			    kp.addr = kallsyms_lookup_name("netif_rx");
			        printk(KERN_EMERG "netif_rx: 0x%x\n", kp.addr);
				    ret = register_kprobe(&kp);
				        if (ret < 0)
						    {
							    	printk(KERN_EMERG "Register kprobe failed\n");
									return -1;
									    }

					    return 0;
}

static void __exit kp_exit(void)
{
	    unregister_kprobe(&kp);
	        printk("Test kp module removed\n");
}

module_init(kp_init);
module_exit(kp_exit);

MODULE_AUTHOR("kkkmmu");
MODULE_DESCRIPTION("Kprobe_Mudule");
MODULE_LICENSE("GPL");

