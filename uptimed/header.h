#include <Cocoa/Cocoa.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <IOKit/IOMessage.h>

// start OSX app
int StartApp();
// set app's menu text
extern void SetLabelText(const char*);
// go func for onSleep
extern void SleepEvent();
// go func for onWake
extern void WakeEvent();
// power event callback
void powerNotificationClbk(void * x, io_service_t y, natural_t messageType, void * messageArgument);
