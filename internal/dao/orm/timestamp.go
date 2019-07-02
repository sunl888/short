package orm

//
//func init() {
//	DefaultCallback.Create().Replace("gorm:create", updateTimeStampForCreateCallback)
//	DefaultCallback.Update().Replace("gorm:update", updateTimeStampForUpdateCallback)
//}
//
//// updateTimeStampForCreateCallback will set `created_at`, `updated_at` when creating
//func updateTimeStampForCreateCallback(scope *Scope) {
//	if !scope.HasError() {
//		now := NowFunc()
//
//		if createdAtField, ok := scope.FieldByName("created_at"); ok {
//			if createdAtField.IsBlank {
//				createdAtField.Set(now)
//			}
//		}
//
//		if updatedAtField, ok := scope.FieldByName("updated_at"); ok {
//			if updatedAtField.IsBlank {
//				updatedAtField.Set(now)
//			}
//		}
//	}
//}
//
//func updateTimeStampForUpdateCallback(scope *Scope) {
//	if _, ok := scope.Get("gorm:update_column"); !ok {
//		scope.SetColumn("updated_at", NowFunc())
//	}
//}
