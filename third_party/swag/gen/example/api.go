package main

type ApiRequest struct {
	// 页码
	Page int `json:"page"`
}

type ApiResponse struct {
	// 列表
	List []string `json:"list"`
	// 总数据
	Total int `json:"total"`
}

type ApiFailResponse struct {
}

// Api
// @Author      CodFrm
// @Summary     一个测试API
// @Description 一个测试API描述
// @ID          example
// @Tags        example
// @Accept      json
// @Produce     json
// @Param       Request body     ApiRequest true "请求信息"
// @Success     200     {object} ApiResponse
// @Failure     400     {object} ApiFailResponse
// @Router      /api [POST]
func Api() {

}
