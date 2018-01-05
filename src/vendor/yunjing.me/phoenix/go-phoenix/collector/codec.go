package collector

import (
	"log"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	pbd "yunjing.me/phoenix/pbd/go/collector"
)

func (self *Collector) encode(zoneId, serviceType, serviceId uint32, happen time.Time, id uint32, args ...interface{}) ([]byte, error) {
	payload := &pbd.Record{
		ZoneId:      zoneId,
		ServiceType: serviceType,
		ServiceId:   serviceId,
		Ts:          happen.Unix(),
		Id:          id,
	}

	if args != nil && len(args) > 0 {
		arr := make([]*pbd.Arg, len(args))
		for k, v := range args {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Bool:
				arr[k] = &pbd.Arg{T: &pbd.Arg_B{B: (v.(bool))}}
			case reflect.Int:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(int))}}
			case reflect.Int8:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(int8))}}
			case reflect.Int16:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(int16))}}
			case reflect.Int32:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: (v.(int32))}}
			case reflect.Int64:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I64{I64: (v.(int64))}}
			case reflect.Uint:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(uint))}}
			case reflect.Uint16:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(uint16))}}
			case reflect.Uint32:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I32{I32: int32(v.(uint32))}}
			case reflect.Uint64:
				arr[k] = &pbd.Arg{T: &pbd.Arg_I64{I64: int64(v.(uint64))}}
			case reflect.Float32:
				arr[k] = &pbd.Arg{T: &pbd.Arg_F32{F32: (v.(float32))}}
			case reflect.Float64:
				arr[k] = &pbd.Arg{T: &pbd.Arg_F64{F64: (v.(float64))}}
			case reflect.String:
				arr[k] = &pbd.Arg{T: &pbd.Arg_S{S: (v.(string))}}
			default:
				log.Printf("数据收集事件%d传递无法识别的数据类型: 第%d字段类型为%s", id, k, reflect.TypeOf(v).String())
				continue
			}
		}
		payload.Args = arr
	}

	return proto.Marshal(payload)
}
