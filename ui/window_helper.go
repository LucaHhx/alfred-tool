package ui

import (
	"runtime"
	"time"
	"unsafe"
	"fyne.io/fyne/v2"
)

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa

#ifdef __APPLE__
#include <Cocoa/Cocoa.h>

// 简化的窗口置顶实现
int setMacWindowAlwaysOnTop(const char* windowTitle, int enable) {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        
        // 遍历所有窗口
        for (NSWindow *window in [app windows]) {
            NSString *title = [window title];
            NSString *targetTitle = [NSString stringWithUTF8String:windowTitle];
            
            // 匹配窗口标题
            if ([title isEqualToString:targetTitle]) {
                dispatch_async(dispatch_get_main_queue(), ^{
                    if (enable) {
                        // 设置窗口层级为最顶层
                        [window setLevel:NSScreenSaverWindowLevel];
                        [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary];
                        [window makeKeyAndOrderFront:nil];
                    } else {
                        // 恢复为正常窗口层级
                        [window setLevel:NSNormalWindowLevel];
                        [window setCollectionBehavior:NSWindowCollectionBehaviorDefault];
                    }
                });
                return 1;
            }
        }
    }
    return 0;
}

// 检查窗口是否置顶
int isMacWindowAlwaysOnTop(const char* windowTitle) {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        
        for (NSWindow *window in [app windows]) {
            NSString *title = [window title];
            NSString *targetTitle = [NSString stringWithUTF8String:windowTitle];
            
            if ([title isEqualToString:targetTitle]) {
                return ([window level] > NSNormalWindowLevel) ? 1 : 0;
            }
        }
    }
    return 0;
}

#else
// 非macOS平台的空实现
int setMacWindowAlwaysOnTop(const char* windowTitle, int enable) {
    return 0;
}

int isMacWindowAlwaysOnTop(const char* windowTitle) {
    return 0;
}
#endif
*/
import "C"

// SetWindowAlwaysOnTop 设置窗口永远置顶
func SetWindowAlwaysOnTop(window fyne.Window, appName string) {
	if runtime.GOOS == "darwin" {
		go func() {
			// 等待窗口完全显示
			time.Sleep(800 * time.Millisecond)
			
			// 获取窗口标题
			windowTitle := window.Title()
			if windowTitle == "" {
				return
			}
			
			// 使用CGO调用设置窗口置顶
			titleCStr := C.CString(windowTitle)
			defer C.free(unsafe.Pointer(titleCStr))
			
			C.setMacWindowAlwaysOnTop(titleCStr, 1)
		}()
	}
}

// IsWindowAlwaysOnTop 检查窗口是否置顶
func IsWindowAlwaysOnTop(window fyne.Window, appName string) bool {
	if runtime.GOOS == "darwin" {
		windowTitle := window.Title()
		if windowTitle == "" {
			return false
		}
		
		titleCStr := C.CString(windowTitle)
		defer C.free(unsafe.Pointer(titleCStr))
		
		result := C.isMacWindowAlwaysOnTop(titleCStr)
		return result == 1
	}
	return false
}

// DisableWindowAlwaysOnTop 取消窗口置顶
func DisableWindowAlwaysOnTop(window fyne.Window, appName string) {
	if runtime.GOOS == "darwin" {
		windowTitle := window.Title()
		if windowTitle == "" {
			return
		}
		
		titleCStr := C.CString(windowTitle)
		defer C.free(unsafe.Pointer(titleCStr))
		
		C.setMacWindowAlwaysOnTop(titleCStr, 0)
	}
}