1. #define container_of(ptr, type, member) ({ \
		     const typeof( ((type *)0)->member ) *__mptr = (ptr); \
			      (type *)( (char *)__mptr - offsetof(type,member) );}) 

2. #define offsetof(type, member) (size_t)&(((type*)0)->member) ------->注意这里是在编译时确定的，所以对于零地址的引用是没有问题的。
