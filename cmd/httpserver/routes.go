package httpserver

// Routes build the routes of the server
func (s *Server) Routes() {
	s.Server.GET("/risk-rules/ping", s.dependencies.StatusHandler.Ping)

	root := s.Server.Group("/risk-rules/v1")

	rulesGroup := root.Group("/rules")
	rulesGroup.POST("", s.dependencies.RulesHandler.AddRule)
	rulesGroup.GET("", s.dependencies.RulesHandler.GetPaged)
	rulesGroup.PUT("/:id", s.dependencies.RulesHandler.UpdateRule)
	rulesGroup.DELETE("/:id", s.dependencies.RulesHandler.RemoveRule)

	chargesGroup := root.Group("/charges")
	chargesGroup.POST("/evaluate", s.dependencies.ChargeHandler.Evaluate)
	chargesGroup.POST("/evaluate_rules", s.dependencies.ChargeHandler.EvaluateOnlyRules)
	chargesGroup.GET("/evaluations/:id", s.dependencies.ChargeHandler.GetEvaluation)
	chargesGroup.GET("/evaluations_rules/:id", s.dependencies.ChargeHandler.GetEvaluationOnlyRules)

	operatorsGroup := root.Group("/operators")
	operatorsGroup.POST("", s.dependencies.OperatorHandler.AddOperator)
	operatorsGroup.GET("", s.dependencies.OperatorHandler.GetAll)
	operatorsGroup.PUT("/:id", s.dependencies.OperatorHandler.Update)
	operatorsGroup.DELETE("/:id", s.dependencies.OperatorHandler.Delete)

	modulesGroup := root.Group("/modules")
	modulesGroup.POST("", s.dependencies.ModulesHandler.Add)
	modulesGroup.GET("", s.dependencies.ModulesHandler.GetAll)
	modulesGroup.DELETE("/:id", s.dependencies.ModulesHandler.Delete)
	modulesGroup.PUT("/:id", s.dependencies.ModulesHandler.Update)

	fieldsGroup := root.Group("/fields")
	fieldsGroup.POST("", s.dependencies.FieldsHandler.Create)
	fieldsGroup.GET("", s.dependencies.FieldsHandler.GetPaged)
	fieldsGroup.DELETE("/:id", s.dependencies.FieldsHandler.Delete)
	fieldsGroup.PUT("/:id", s.dependencies.FieldsHandler.Update)

	conditionsGroup := root.Group("/conditions")
	conditionsGroup.POST("", s.dependencies.ConditionsHandler.Add)
	conditionsGroup.GET("", s.dependencies.ConditionsHandler.GetPaged)
	conditionsGroup.PUT("/:id", s.dependencies.ConditionsHandler.Update)
	conditionsGroup.DELETE("/:id", s.dependencies.ConditionsHandler.Delete)

	familiesGroup := root.Group("/families")
	familiesGroup.POST("", s.dependencies.FamilyHandler.Create)
	familiesGroup.DELETE("/:id", s.dependencies.FamilyHandler.Delete)
	familiesGroup.PUT("/:id", s.dependencies.FamilyHandler.Update)
	familiesGroup.GET("", s.dependencies.FamilyHandler.Get)

	familyCompaniesGroup := root.Group("/family_companies")
	familyCompaniesGroup.POST("", s.dependencies.FamilyCompaniesHandler.Create)
	familyCompaniesGroup.PUT("/:id", s.dependencies.FamilyCompaniesHandler.Update)
	familyCompaniesGroup.DELETE("/:id", s.dependencies.FamilyCompaniesHandler.Delete)
	familyCompaniesGroup.GET("", s.dependencies.FamilyCompaniesHandler.Get)

	merchantsGroup := root.Group("/merchants_score")
	merchantsGroup.POST("", s.dependencies.MerchantsScoreHandler.MerchantScoreProcessing)
}
