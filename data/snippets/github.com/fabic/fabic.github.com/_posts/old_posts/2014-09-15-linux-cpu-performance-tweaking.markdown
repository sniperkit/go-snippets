---
layout: post
title: "Linux CPU frequency troobleshooting"
description: "Sequel of yesterday's post, CPU was in powersave mode."
category: diary
tags: [Linux, CPU, troobleshooting, old]
---

_Sequel of yesterday's post about compiling LLVM Clang: it hadn't completed this morning due to my laptop's CPU "governor" thing being "powersave"._

## The morning after -- Sept. 15th -- 50% done only ?! o\_o\`

_It shall have completed a long while ago -- 15 min. load avg. is 3.6 -- compilation is still goin' on with 3 jobs using >80% of CPUs -- what the heck ? -- damn laptop must have gone in low-power consumption mode somehow -- I do remember having closed the lid of it inadvertently before dozzing off: it went in sleep mode, I woke it up and stuff continued ; but probably did that trigger some sort of CPU powersave state._

### Monitoring CPU frequency :

Shows that it hardly ever gets above 1GHz frequency :

	[fabi@wall] ~/dev/llvm/build (master ✗)  watch -d -n1 \
	    grep -E "'^(processor|cpu MHz|model name|$)'" /proc/cpuinfo

	processor	: 0
	model name	: Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz
	cpu MHz		: 427.781

	processor	: 1
	model name	: Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz
	cpu MHz		: 411.750

	processor	: 2
	model name	: Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz
	cpu MHz		: 427.675

	processor	: 3
	model name	: Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz
	cpu MHz		: 410.800

### Performance troobleshooting :

**Did try** to unplug-plug the power cable, didn't do it ;

**Did try** to start `laptop-mode` daemon too, didn't do it -- stopped it as it tends to spin down the hard-drive too often ;

Reading [archlinux's CPU frequency scaling](https://wiki.archlinux.org/index.php/CPU_Frequency_Scaling#Scaling_governors), and investigating how to trigger usage of a _performance_ CPU frequency governor :

	[root@wall] lib # find /proc /sys -name \*pstate\*
	/sys/devices/system/cpu/intel_pstate
	/sys/kernel/debug/pstate_snb
	/sys/module/acpi_cpufreq/parameters/acpi_pstate_strict

Seems not to be under `/sys/devices/system/cpu/intel_pstate/` :

	[root@wall] lib # grep ^ /sys/devices/system/cpu/intel_pstate/* | column -t -s:
	/sys/devices/system/cpu/intel_pstate/max_perf_pct  100
	/sys/devices/system/cpu/intel_pstate/min_perf_pct  25
	/sys/devices/system/cpu/intel_pstate/no_turbo      0


	[root@wall] lib # uname -a
	Linux sabayon 3.12.0-sabayon #1 SMP Mon Sep 1 02:30:34 UTC 2014 x86_64 Intel(R) Core(TM) i7-2620M CPU @ 2.70GHz GenuineIntel GNU/Linux


	[root@wall] lib # ls /lib/modules/$(uname -r)/kernel/drivers/cpufreq/
	amd_freq_sensitivity.ko  pcc-cpufreq.ko

##### **Found it**, we're indeed in powersave mode :

	[root@wall] lib # find /proc /sys -name \*cpufreq\*
	/sys/devices/system/cpu/cpu0/cpufreq
	/sys/devices/system/cpu/cpu1/cpufreq
	/sys/devices/system/cpu/cpu2/cpufreq
	/sys/devices/system/cpu/cpu3/cpufreq
	/sys/module/acpi_cpufreq

	[root@wall] lib # grep ^ /sys/devices/system/cpu/cpu?/cpufreq/scaling_governor
	/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor:powersave
	/sys/devices/system/cpu/cpu1/cpufreq/scaling_governor:powersave
	/sys/devices/system/cpu/cpu2/cpufreq/scaling_governor:powersave
	/sys/devices/system/cpu/cpu3/cpufreq/scaling_governor:powersave

##### Availabe CPU frequency governors :

	[root@wall] lib # grep ^ /sys/devices/system/cpu/cpu?/cpufreq/scaling_available_governors
	/sys/devices/system/cpu/cpu0/cpufreq/scaling_available_governors:performance powersave
	/sys/devices/system/cpu/cpu1/cpufreq/scaling_available_governors:performance powersave
	/sys/devices/system/cpu/cpu2/cpufreq/scaling_available_governors:performance powersave
	/sys/devices/system/cpu/cpu3/cpufreq/scaling_available_governors:performance powersave

##### **Ok**, this does affect CPU frequency :

	[root@wall] lib # for cpu in /sys/devices/system/cpu/cpu?/cpufreq/scaling_governor ; do echo powersave > $cpu; done

	[root@wall] lib # for cpu in /sys/devices/system/cpu/cpu?/cpufreq/scaling_governor ; do echo performance > $cpu; done

### Searching for more powa!

**Fine**, CPUs remain at 1.7 GHz under load, with seldom peaks at 2 GHz -- pursuing investigations searching for "more performance" or reactivity :

	[root@wall] lib # for thing in /sys/devices/system/cpu/cpu?/cpufreq/* ; do \
	      echo `echo $thing|grep -o 'cpu[0-9]'`:`basename $thing`:`cat $thing`:--:$thing ;\
	    done |
	      column -s: -t

	cpu0  affected_cpus                0
	cpu0  cpuinfo_cur_freq             1699945
	cpu0  cpuinfo_max_freq             3400000
	cpu0  cpuinfo_min_freq             800000
	cpu0  cpuinfo_transition_latency   4294967295
	cpu0  related_cpus                 0
	cpu0  scaling_available_governors  performance powersave
	cpu0  scaling_driver               intel_pstate
	cpu0  scaling_governor             performance
	cpu0  scaling_max_freq             3400000
	cpu0  scaling_min_freq             800000
	cpu0  scaling_setspeed             <unsupported>

_Hum... -- reading [archlinux's Maximizing performance](https://wiki.archlinux.org/index.php/Maximizing_Performance)..._

##### Hard-drive read speed seems **ok** :

	[root@wall] lib # hdparm -t /dev/sda

	/dev/sda:
	 Timing buffered disk reads: 298 MB in  3.01 seconds =  99.01 MB/sec

##### Free memory (with 3 Clang jobs running) seems **ok** :

	[root@wall] lib # free -h
	             total       used       free     shared    buffers     cached
	Mem:          3.7G       3.3G       447M       133M        75M       1.3G
	-/+ buffers/cache:       1.9G       1.8G
	Swap:         4.0G        18M       4.0G

##### Reading [Linux kernel doc.: cpu-freq/governors.txt](https://www.kernel.org/doc/Documentation/cpu-freq/governors.txt)

Can't set the `cpuinfo_transition_latency` :-/

	[root@wall] lib # echo 594967295 > /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_transition_latency
-su: /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_transition_latency: Permission denied


_LLVM build is at 64% -- gonna have some coffee -- will read some doc. probably..._

_**End of transmission**, the morning after, Sept. 15th._
