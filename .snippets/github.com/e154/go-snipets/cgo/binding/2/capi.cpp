
#include "capi.h"

SystemTray_ NewSystemTray() {
	SystemTray * t = new SystemTray();
	return (void*)t;
}

void SetTrayIcon(SystemTray_ t, char *img) {
	SystemTray * ptr = (SystemTray*)t;
	ptr->setTrayIcon(img);
}

void SetTrayToolTip(SystemTray_ t, char *tooltip) {
	SystemTray * ptr = (SystemTray*)t;
	ptr->setTrayToolTip(tooltip);
}

void Init(SystemTray_ t) {
	SystemTray * ptr = (SystemTray*)t;
	ptr->init();
}