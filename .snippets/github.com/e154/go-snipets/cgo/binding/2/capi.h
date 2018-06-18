/** Copyright (C), DeltaSync Studios, 2010-2014. All rights reserved.
 * ------------------------------------------------------------------
 * File name:   capi.h
 * Version:     v1.00
 * Created:     23:42:52 / 18 авг. 2015 г.
 * Author:      delta54 <support@e154.ru>
 *
 * Your use and or redistribution of this software in source and / or
 * binary form, with or without modification, is subject to: (i) your
 * ongoing acceptance of and compliance with the terms and conditions of
 * the DeltaSync License Agreement; and (ii) your inclusion of this notice
 * in any version of this software that you use or redistribute.
 * A copy of the DeltaSync License Agreement is available by contacting
 * DeltaSync Studios. at http://www.inet-print.ru/
 *
 * Description:
 * ------------------------------------------------------------------
 * History:
 *
 */

#ifndef CAPI_CAPI_H_
#define CAPI_CAPI_H_

#ifdef __cplusplus
extern "C" {
#endif

// SystemTray
typedef void* SystemTray_;
SystemTray_ NewSystemTray(void);
void SetTrayIcon(SystemTray_, char*);
void SetTrayToolTip(SystemTray_, char*);
void Init(SystemTray_);

#ifdef __cplusplus
} // extern "C"
#endif

#endif /* CAPI_CAPI_H_ */
