package appmock

//go:generate moq -out gen_auth_store_moq.go -pkg appmock ../app AuthStore
//go:generate moq -out gen_auth_respository_moq.go -pkg appmock ../app AuthRepository
//go:generate moq -out gen_token_issuer_moq.go -pkg appmock ../app TokenIssuer
//go:generate moq -out gen_phc_builder_factory_moq.go -pkg appmock ../app PHCBuilderFactory
//go:generate moq -out gen_external_identity_resolver_moq.go -pkg appmock ../app ExternalIdentityResolver
//go:generate moq -out gen_external_identity_resolver_registry_moq.go -pkg appmock ../app ExternalIdentityResolverRegistry
