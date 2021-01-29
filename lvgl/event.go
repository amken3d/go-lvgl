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
  extern void go_event_callback(lv_obj_t * obj, lv_event_t event);

  // callback 'proxy'
  static inline void event_cb_proxy(lv_obj_t * obj, lv_event_t event) {
	go_event_callback(obj, event);
  }

  // Callback registration
  static void _register_callback(lv_obj_t *obj) {
	  lv_obj_set_event_cb(obj, event_cb_proxy);
  }

*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
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
	IDX int
}

// EventCallbackFn is the callback function prototype
type EventCallbackFn func(*LVObj, LVEvent)

// LVEvent represents lv_event_t
type LVEvent C.lv_event_t

// EventFnRegister contains a map of all event callbacks
type eventFnRegister struct {
	sync.Mutex
	index int
	fns   map[int]EventCallbackFn
}

var (
	callbackRegister eventFnRegister
)

func init() {
	callbackRegister.fns = make(map[int]EventCallbackFn)
}

func (r eventFnRegister) register(fn EventCallbackFn) int {
	r.Lock()
	defer r.Unlock()
	r.index++
	for r.fns[callbackRegister.index] != nil {
		r.index++
	}
	r.fns[callbackRegister.index] = fn
	return r.index
}

func (r eventFnRegister) lookup(i int) EventCallbackFn {
	r.Lock()
	defer r.Unlock()
	return r.fns[i]
}

func (r eventFnRegister) unregister(i int) {
	r.Lock()
	defer r.Unlock()
	delete(r.fns, i)
}

//export go_event_callback
func go_event_callback(obj *C.struct__lv_obj_t, event C.lv_event_t) {
	fmt.Printf("In Callback: fns = %v+\n", callbackRegister.fns)
	o := (*LVObj)(obj)
	eud := pointer.Restore(unsafe.Pointer(o.UserData())).(*EventUserData)
	fn := callbackRegister.lookup(eud.IDX)
	fn(o, (LVEvent)(event))
}

// MyCallback is a test callback function
func MyCallback(obj *LVObj, event LVEvent) {
	fmt.Printf("Received event: %v\n", event)
}

// TryCallback is a test function
func (obj *LVObj) TryCallback() {
	i := callbackRegister.register(MyCallback)
	obj.SetUserData(pointer.Save(&EventUserData{IDX: i}))
	C._register_callback(((*C.struct__lv_obj_t)(unsafe.Pointer(obj))))
	fmt.Printf("After register: fns = %v+\n", callbackRegister.fns)
}
