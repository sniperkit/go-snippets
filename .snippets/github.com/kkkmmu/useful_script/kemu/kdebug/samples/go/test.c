Tracing "ip_rcv"... Ctrl-C to end.
 0)               |  ip_rcv() {
 0)               |    nf_hook_slow() {
 0)               |      nf_iterate() {
 0)   7.067 us    |        ipv4_conntrack_defrag();
 0)               |        ipv4_conntrack_in() {
 0)               |          nf_conntrack_in() {
 0)   5.563 us    |            ipv4_get_l4proto();
 0)   4.398 us    |            __nf_ct_l4proto_find();
 0)               |            udp_error() {
 0)   3.762 us    |              nf_ip_checksum();
 0) + 31.426 us   |            }
 0)               |            nf_ct_get_tuple() {
 0)   4.380 us    |              ipv4_pkt_to_tuple();
 0)   3.584 us    |              udp_pkt_to_tuple();
 0) + 53.434 us   |            }
 0)               |            __nf_conntrack_find_get() {
 0)   5.388 us    |              __local_bh_enable_ip();
 0) + 27.306 us   |            }
 0)               |            init_conntrack() {
 0)               |              nf_ct_invert_tuple() {
 0)   2.889 us    |                ipv4_invert_tuple();
 0)   2.913 us    |                udp_invert_tuple();
 0) + 48.546 us   |              }
 0)               |              __nf_conntrack_alloc.isra.42() {
 0)   6.034 us    |                kmem_cache_alloc();
 0)   3.123 us    |                init_timer_key();
 0) + 50.469 us   |              }
 0)   2.931 us    |              udp_get_timeouts();
 0)   2.767 us    |              udp_new();
 0)   6.795 us    |              __nf_ct_try_assign_helper();
 0)   3.718 us    |              _raw_spin_lock();
 0)   3.328 us    |              __local_bh_enable_ip();
 0) ! 255.833 us  |            }
 0)   2.657 us    |            udp_get_timeouts();
 0)               |            udp_packet() {
 0)   4.152 us    |              __nf_ct_refresh_acct();
 0) + 28.192 us   |            }
 0) ! 581.752 us  |          }
 0) ! 607.076 us  |        }
 0)               |        iptable_mangle_hook() {
 0)               |          ipt_do_table() {
 0)   3.071 us    |            __local_bh_enable_ip();
 0) + 28.758 us   |          }
 0) + 54.630 us   |        }
 0) ! 743.980 us  |      }
 0) ! 768.572 us  |    }
 0)               |    ip_rcv_finish() {
 0)   7.545 us    |      udp_v4_early_demux();
 0)               |      ip_route_input_noref() {
 0)               |        ip_route_input_slow() {
 0) + 12.507 us   |          fib_table_lookup();
 0)               |          fib_validate_source() {
 0)   6.946 us    |            fib_table_lookup();
 0) + 29.777 us   |          }
 0) + 88.287 us   |        }
 0) ! 108.666 us  |      }
 0)               |      ip_local_deliver() {
 0)               |        nf_hook_slow() {
 0)               |          nf_iterate() {
 0)               |            iptable_mangle_hook() {
 0)               |              ipt_do_table() {
 0)   3.063 us    |                __local_bh_enable_ip();
 0) + 24.014 us   |              }
 0) + 42.402 us   |            }
 0)               |            iptable_filter_hook() {
 0)               |              ipt_do_table() {
 0)   2.999 us    |                __local_bh_enable_ip();
 0) + 23.490 us   |              }
 0) + 43.998 us   |            }
 0)   3.617 us    |            ipv4_helper();
 0)               |            ipv4_confirm() {
 0)               |              __nf_conntrack_confirm() {
 0)   3.863 us    |                __hash_conntrack();
 0)               |                nf_conntrack_double_lock() {
 0)               |                  nf_conntrack_lock() {
 0)   3.754 us    |                    _raw_spin_lock();
 0) + 23.518 us   |                  }
 0)   3.490 us    |                  _raw_spin_lock();
 0) + 74.507 us   |                }
 0)               |                nf_ct_del_from_dying_or_unconfirmed_list() {
 0)   3.162 us    |                  _raw_spin_lock();
 0) + 23.753 us   |                }
 0)               |                add_timer() {
 0)               |                  mod_timer() {
 0)               |                    lock_timer_base.isra.22() {
 0)   3.328 us    |                      _raw_spin_lock_irqsave();
 0) + 22.121 us   |                    }/*asdkfjdsakfjsadkfjkdf*/
 0)   2.750 us    |                    detach_if_pending();
 0)               |                    get_nohz_timer_target() {
 0)   3.124 us    |                      idle_cpu();
 0) + 21.215 us   |                    }
 0)               |                    internal_add_timer() {
 0)   3.684 us    |                      __internal_add_timer();
 0)   2.941 us    |                      wake_up_nohz_cpu();
 0) + 40.135 us   |                    }
 0)   2.825 us    |                    _raw_spin_unlock_irqrestore();
 0) ! 169.596 us  |                  }
 0) ! 188.980 us  |                }
 0)   3.757 us    |                __nf_conntrack_hash_insert();
 0)   3.657 us    |                nf_conntrack_double_unlock();
 0)   3.119 us    |                __local_bh_enable_ip();
 0) ! 439.633 us  |              }
 0) ! 466.094 us  |            }
 0) ! 641.685 us  |          }
 0) ! 659.424 us  |        }
 0)               |        ip_local_deliver_finish() {
 0)   4.294 us    |          raw_local_deliver();
 0)               |          udp_rcv() {
 0)               |            __udp4_lib_rcv() {
 0)               |              __udp4_lib_mcast_deliver() {
 0)   3.206 us    |                _raw_spin_lock();
 0)               |                consume_skb() {
 0)               |                  skb_release_all() {
 0)   4.889 us    |                    skb_release_head_state();
 0)               |                    skb_release_data() {
 0)               |                      skb_free_head() {
 0)   6.541 us    |                        __free_page_frag();
 0) + 26.581 us   |                      }
 0) + 47.177 us   |                    }
 0) + 88.658 us   |                  }
 0)               |                  kfree_skbmem() {
 0)   8.359 us    |                    kmem_cache_free();
 0) + 29.367 us   |                  }
 0) ! 153.623 us  |                }
 0) ! 195.836 us  |              }
 0) ! 225.281 us  |            }
 0) ! 248.067 us  |          }
 0) ! 299.641 us  |        }
 0) ! 999.333 us  |      }
 0) # 1183.036 us |    }
 0) # 2018.213 us |  }

Ending tracing...
