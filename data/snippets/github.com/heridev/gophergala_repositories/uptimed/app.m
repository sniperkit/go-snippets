#include "header.h"

io_connect_t root_port;

// register onWake and onSleep events
void registerPowerNotificationCallbacks() {
  IONotificationPortRef notifyPortRef;
  io_object_t notifierObject;
  root_port = IORegisterForSystemPower(&root_port, &notifyPortRef, powerNotificationClbk, &notifierObject);
  if (!root_port) 
    exit (1);
  CFRunLoopAddSource (CFRunLoopGetCurrent(), IONotificationPortGetRunLoopSource(notifyPortRef), kCFRunLoopDefaultMode);
}

// forward events to go via SleepEvent and WakeEvent
void powerNotificationClbk(void * x, io_service_t y, natural_t messageType, void * messageArgument) {
  switch (messageType) {
    case kIOMessageCanSystemSleep:
      IOAllowPowerChange(root_port, (long)messageArgument);
      break;

    case kIOMessageSystemWillSleep:
      IOAllowPowerChange(root_port, (long)messageArgument);
      SleepEvent();
      break;

    case kIOMessageSystemHasPoweredOn:
      //System has finished waking up...
      WakeEvent();
      break;

    default:
      break;
  }
}

NSStatusItem *statusItem;

int StartApp(void) {
  [NSAutoreleasePool new];
  [NSApplication sharedApplication];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited]; // hide from dock

  registerPowerNotificationCallbacks();

  NSMenuItem *tItem = nil;
  NSMenu *menu;

  menu = [[NSMenu alloc] initWithTitle:@""];
  [menu setAutoenablesItems:NO];
  tItem = [menu addItemWithTitle:@"Quit" action:@selector(terminate:) keyEquivalent:@"q"];
  [tItem setKeyEquivalentModifierMask:NSCommandKeyMask];

  NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
  statusItem = [statusBar statusItemWithLength:NSVariableStatusItemLength];
  [statusItem retain];
  //[statusItem setTitle:@"00:00:00"];
  [statusItem setHighlightMode:YES];
  [statusItem setMenu:menu];

  [NSApp activateIgnoringOtherApps:YES];
  [NSApp run];
  return 0;
}

// Set the app's menu label. Called from go
void SetLabelText(const char *str) {
  @autoreleasepool {
    // TODO: nill check on statusItem
    NSString *text = [NSString stringWithUTF8String:str];
    [statusItem setTitle:text];
  }
}
