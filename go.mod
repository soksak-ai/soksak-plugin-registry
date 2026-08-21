module github.com/soksak-ai/soksak-plugin-registry

go 1.25.0

require github.com/soksak-ai/soksak-contract-registry v0.0.0-20260821065357-22e40ef8d315

require github.com/soksak-ai/soksak-contract-composition v0.0.0-20260821062043-838e68be8996 // indirect

replace github.com/soksak-ai/soksak-contract-composition => ../soksak-contracts/soksak-contract-composition

replace github.com/soksak-ai/soksak-contract-registry => ../soksak-contracts/soksak-contract-registry
