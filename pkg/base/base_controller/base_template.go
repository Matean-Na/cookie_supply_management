package base_controller

import (
	"cookie_supply_management/pkg/base/base_service"
	"cookie_supply_management/utils/parsers"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type CrudTemplateInterface interface {
	FindAll(ctx *gin.Context) *AppError
	FindOne(ctx *gin.Context) *AppError
	Create(ctx *gin.Context) *AppError
	Update(ctx *gin.Context) *AppError
	Delete(ctx *gin.Context) *AppError
}

type CrudTemplate struct {
	mi ModelInterface
	ri RedisInterface
}

func NewCrudTemplate(cc *CrudController) *CrudTemplate {
	ct := &CrudTemplate{
		mi: cc.ModelInterface,
		ri: cc,
	}
	return ct
}

func (ct *CrudTemplate) FindAll(ctx *gin.Context, allInter base_service.FindAllInterface) *AppError {
	return ct.FindAllFunc(ctx, allInter.FindAll)
}

func (ct *CrudTemplate) FindAllFunc(ctx *gin.Context, findAll base_service.FindAll) *AppError {
	model := ct.mi.GetAll()
	pager := GetPager(ctx)

	var total int64
	if err := findAll(&model, ct.mi.ScopeAll, pager, GetOrder(ctx, ct.mi.GetOne()), &total, GetQuery(ctx, model), ct.ri.KeyAll()); err != nil {
		return I18nError(ctx, model, "exception:could-not-fetch-records")
	}

	return OkT(ctx, total, model)
}

func (ct *CrudTemplate) FindOne(ctx *gin.Context, oneInter base_service.FindOneInterface) *AppError {
	return ct.FindOneFunc(ctx, oneInter.FindOne)
}

func (ct *CrudTemplate) FindOneFunc(ctx *gin.Context, findOne base_service.FindOne) *AppError {
	id := ctx.Param("id")
	model := ct.mi.GetOne()
	model.SetId(parsers.ParamUint(id))

	if err := findOne(model, ct.mi.ScopeOne); err != nil {
		return ErrNotFound(ctx, model, model.GetId(), "id")
	}

	return Ok(ctx, model)
}

func (ct *CrudTemplate) Create(ctx *gin.Context, creInter base_service.CreateInterface) *AppError {
	return ct.CreateFunc(ctx, creInter.Create)
}

func (ct *CrudTemplate) CreateFunc(ctx *gin.Context, create base_service.Create) *AppError {
	i := ct.mi.GetOne()

	if err := ctx.ShouldBindJSON(i); err != nil {
		return BaseError(err.Error(), http.StatusBadRequest)
	}

	//set code in model from context
	lanCode := parsers.ParamUint(ctx.GetString("language_code"))
	sType := reflect.ValueOf(i).Elem()
	field := sType.FieldByName("LanguageID")
	if field.IsValid() {
		field.SetUint(uint64(lanCode))
	}

	if err := create(i, ct.ri.KeyAll()); err != nil {
		return ErrNotCreated(err)
	}

	return Ok(ctx, gin.H{"id": i.GetId()})
}

func (ct *CrudTemplate) UpdateFunc(ctx *gin.Context, update base_service.Update) *AppError {
	id := ctx.Param("id")
	model := ct.mi.GetOne()

	if err := ctx.ShouldBindJSON(model); err != nil {
		return LocalizeError(ctx, err)
	}

	model.SetId(parsers.ParamUint(id))

	if err := update(model, ct.ri.KeyAll()); err != nil {
		return ErrNotUpdated(err)
	}

	return Ok(ctx, gin.H{"id": model.GetId()})
}

func (ct *CrudTemplate) Update(ctx *gin.Context, updInter base_service.UpdateInterface) *AppError {
	return ct.UpdateFunc(ctx, updInter.Update)
}

func (ct *CrudTemplate) DeleteFunc(ctx *gin.Context, delete base_service.Delete) *AppError {
	id := ctx.Param("id")

	model := ct.mi.GetOne()
	model.SetId(parsers.ParamUint(id))

	if err := delete(model, ct.ri.KeyAll()); err != nil {
		return ErrNotDeleted(err)
	}

	return Ok(ctx, gin.H{"id": model.GetId()})
}

func (ct *CrudTemplate) Delete(ctx *gin.Context, delInter base_service.DeleteInterface) *AppError {
	return ct.DeleteFunc(ctx, delInter.Delete)
}
