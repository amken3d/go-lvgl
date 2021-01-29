package lvgl

/*
  #cgo CFLAGS: -I../include -I../include/lvgl -I../include/lv_drivers/display -I../include/lv_drivers/indev
  #include "lv_conf.h"
  #include "lv_drv_conf.h"
  #include "lvgl.h"
  #include "fbdev.h"
  #include "evdev.h"

  #cgo LDFLAGS: -L../include -llvgl

  // Go function prototype
  extern void go_callback_int(int foo, int p1);

  // callback 'proxy'
  static inline void CallMyFunction(int foo) {
	go_callback_int(foo, 5);
  }

*/
import "C"
import (
	"fmt"
	"sync"
)

const (
	LVEventPressed           uint8 = C.LV_EVENT_PRESSED
	LVEventPressing          uint8 = C.LV_EVENT_PRESSING
	LVEventPressLost         uint8 = C.LV_EVENT_PRESS_LOST
	LVEventShortClicked      uint8 = C.LV_EVENT_SHORT_CLICKED
	LVEventLongPressed       uint8 = C.LV_EVENT_LONG_PRESSED
	LVEventLongPressedRepeat uint8 = C.LV_EVENT_LONG_PRESSED_REPEAT
	LVEventClicked           uint8 = C.LV_EVENT_CLICKED
	LVEventReleased          uint8 = C.LV_EVENT_RELEASED
	LVEventDragBegin         uint8 = C.LV_EVENT_DRAG_BEGIN
	LVEventDragEnd           uint8 = C.LV_EVENT_DRAG_END
	LVEventDragThrowBegin    uint8 = C.LV_EVENT_DRAG_THROW_BEGIN
	LVEventGesture           uint8 = C.LV_EVENT_GESTURE
	LVEventKey               uint8 = C.LV_EVENT_KEY
	LVEventFocused           uint8 = C.LV_EVENT_FOCUSED
	LVEventDefocused         uint8 = C.LV_EVENT_DEFOCUSED
	LVEventLeave             uint8 = C.LV_EVENT_LEAVE
	LVEventValueChanged      uint8 = C.LV_EVENT_VALUE_CHANGED
	LVEventInsert            uint8 = C.LV_EVENT_INSERT
	LVEventRefresh           uint8 = C.LV_EVENT_REFRESH
	LVEventApply             uint8 = C.LV_EVENT_APPLY
	LVEventCancel            uint8 = C.LV_EVENT_CANCEL
	LVEventDelete            uint8 = C.LV_EVENT_DELETE
)

// EventUserData represents event info to inject into lv_obj as user_data
type EventUserData struct {
	idx int
}

// EventCallbackFn is the callback function prototype
type EventCallbackFn func(*LVObj, LVEvent)

// LVEvent represents lv_event_t
type LVEvent C.lv_event_t

// Callback contains a datastructure to share with C code
type Callback struct {
	Func   EventCallbackFn
	Object *LVObj
}

var mu sync.Mutex
var index int
var fns = make(map[int]func(C.int))

func register(fn func(C.int)) int {
	mu.Lock()
	defer mu.Unlock()
	index++
	for fns[index] != nil {
		index++
	}
	fns[index] = fn
	return index
}

func lookup(i int) func(C.int) {
	mu.Lock()
	defer mu.Unlock()
	return fns[i]
}

func unregister(i int) {
	mu.Lock()
	defer mu.Unlock()
	delete(fns, i)
}

//export go_callback_int
func go_callback_int(foo C.int, p1 C.int) {
	fn := lookup(int(foo))
	fn(p1)
}

func MyCallback(x C.int) {
	fmt.Println("callback with", x)
}

func (obj *LVObj) TryCallback() {
	i := register(MyCallback)
	C.CallMyFunction(C.int(i))
	unregister(i)
}
