
#include "capi.h"

SystemTray_ *NewSystemTray() { return new SystemTray(); }
void SetTrayIcon(SystemTray_ *t, char *img) { static_cast<SystemTray *>(t)->setTrayIcon(img); }
void SetTrayToolTip(SystemTray_ *t, char *tooltip) { static_cast<SystemTray *>(t)->setTrayToolTip(tooltip); }
void Init(SystemTray_ *t) { static_cast<SystemTray *>(t)->init(); }