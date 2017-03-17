package pbgeo

import (
	"bytes"
	"fmt"
	"strings"

	"go.pedge.io/pb/gogo/pb/money"
)

var (
	//CountryCodeAD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AD for simplicity.
	CountryCodeAD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AD
	//CountryCodeAE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AE for simplicity.
	CountryCodeAE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AE
	//CountryCodeAF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AF for simplicity.
	CountryCodeAF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AF
	//CountryCodeAG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AG for simplicity.
	CountryCodeAG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AG
	//CountryCodeAI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AI for simplicity.
	CountryCodeAI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AI
	//CountryCodeAL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AL for simplicity.
	CountryCodeAL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AL
	//CountryCodeAM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AM for simplicity.
	CountryCodeAM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AM
	//CountryCodeAO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AO for simplicity.
	CountryCodeAO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AO
	//CountryCodeAQ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AQ for simplicity.
	CountryCodeAQ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AQ
	//CountryCodeAR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AR for simplicity.
	CountryCodeAR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AR
	//CountryCodeAS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AS for simplicity.
	CountryCodeAS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AS
	//CountryCodeAT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AT for simplicity.
	CountryCodeAT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AT
	//CountryCodeAU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AU for simplicity.
	CountryCodeAU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AU
	//CountryCodeAW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AW for simplicity.
	CountryCodeAW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AW
	//CountryCodeAX is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AX for simplicity.
	CountryCodeAX = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AX
	//CountryCodeAZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AZ for simplicity.
	CountryCodeAZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_AZ
	//CountryCodeBA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BA for simplicity.
	CountryCodeBA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BA
	//CountryCodeBB is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BB for simplicity.
	CountryCodeBB = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BB
	//CountryCodeBD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BD for simplicity.
	CountryCodeBD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BD
	//CountryCodeBE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BE for simplicity.
	CountryCodeBE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BE
	//CountryCodeBF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BF for simplicity.
	CountryCodeBF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BF
	//CountryCodeBG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BG for simplicity.
	CountryCodeBG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BG
	//CountryCodeBH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BH for simplicity.
	CountryCodeBH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BH
	//CountryCodeBI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BI for simplicity.
	CountryCodeBI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BI
	//CountryCodeBJ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BJ for simplicity.
	CountryCodeBJ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BJ
	//CountryCodeBL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BL for simplicity.
	CountryCodeBL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BL
	//CountryCodeBM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BM for simplicity.
	CountryCodeBM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BM
	//CountryCodeBN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BN for simplicity.
	CountryCodeBN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BN
	//CountryCodeBO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BO for simplicity.
	CountryCodeBO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BO
	//CountryCodeBQ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BQ for simplicity.
	CountryCodeBQ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BQ
	//CountryCodeBR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BR for simplicity.
	CountryCodeBR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BR
	//CountryCodeBS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BS for simplicity.
	CountryCodeBS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BS
	//CountryCodeBT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BT for simplicity.
	CountryCodeBT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BT
	//CountryCodeBV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BV for simplicity.
	CountryCodeBV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BV
	//CountryCodeBW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BW for simplicity.
	CountryCodeBW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BW
	//CountryCodeBY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BY for simplicity.
	CountryCodeBY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BY
	//CountryCodeBZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BZ for simplicity.
	CountryCodeBZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_BZ
	//CountryCodeCA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CA for simplicity.
	CountryCodeCA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CA
	//CountryCodeCC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CC for simplicity.
	CountryCodeCC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CC
	//CountryCodeCD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CD for simplicity.
	CountryCodeCD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CD
	//CountryCodeCF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CF for simplicity.
	CountryCodeCF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CF
	//CountryCodeCG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CG for simplicity.
	CountryCodeCG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CG
	//CountryCodeCH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CH for simplicity.
	CountryCodeCH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CH
	//CountryCodeCI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CI for simplicity.
	CountryCodeCI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CI
	//CountryCodeCK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CK for simplicity.
	CountryCodeCK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CK
	//CountryCodeCL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CL for simplicity.
	CountryCodeCL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CL
	//CountryCodeCM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CM for simplicity.
	CountryCodeCM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CM
	//CountryCodeCN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CN for simplicity.
	CountryCodeCN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CN
	//CountryCodeCO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CO for simplicity.
	CountryCodeCO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CO
	//CountryCodeCR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CR for simplicity.
	CountryCodeCR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CR
	//CountryCodeCU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CU for simplicity.
	CountryCodeCU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CU
	//CountryCodeCV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CV for simplicity.
	CountryCodeCV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CV
	//CountryCodeCW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CW for simplicity.
	CountryCodeCW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CW
	//CountryCodeCX is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CX for simplicity.
	CountryCodeCX = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CX
	//CountryCodeCY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CY for simplicity.
	CountryCodeCY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CY
	//CountryCodeCZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CZ for simplicity.
	CountryCodeCZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_CZ
	//CountryCodeDE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DE for simplicity.
	CountryCodeDE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DE
	//CountryCodeDJ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DJ for simplicity.
	CountryCodeDJ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DJ
	//CountryCodeDK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DK for simplicity.
	CountryCodeDK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DK
	//CountryCodeDM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DM for simplicity.
	CountryCodeDM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DM
	//CountryCodeDO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DO for simplicity.
	CountryCodeDO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DO
	//CountryCodeDZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DZ for simplicity.
	CountryCodeDZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_DZ
	//CountryCodeEC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EC for simplicity.
	CountryCodeEC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EC
	//CountryCodeEE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EE for simplicity.
	CountryCodeEE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EE
	//CountryCodeEG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EG for simplicity.
	CountryCodeEG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EG
	//CountryCodeEH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EH for simplicity.
	CountryCodeEH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_EH
	//CountryCodeER is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ER for simplicity.
	CountryCodeER = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ER
	//CountryCodeES is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ES for simplicity.
	CountryCodeES = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ES
	//CountryCodeET is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ET for simplicity.
	CountryCodeET = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ET
	//CountryCodeFI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FI for simplicity.
	CountryCodeFI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FI
	//CountryCodeFJ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FJ for simplicity.
	CountryCodeFJ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FJ
	//CountryCodeFK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FK for simplicity.
	CountryCodeFK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FK
	//CountryCodeFM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FM for simplicity.
	CountryCodeFM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FM
	//CountryCodeFO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FO for simplicity.
	CountryCodeFO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FO
	//CountryCodeFR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FR for simplicity.
	CountryCodeFR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_FR
	//CountryCodeGA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GA for simplicity.
	CountryCodeGA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GA
	//CountryCodeGB is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GB for simplicity.
	CountryCodeGB = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GB
	//CountryCodeGD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GD for simplicity.
	CountryCodeGD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GD
	//CountryCodeGE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GE for simplicity.
	CountryCodeGE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GE
	//CountryCodeGF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GF for simplicity.
	CountryCodeGF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GF
	//CountryCodeGG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GG for simplicity.
	CountryCodeGG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GG
	//CountryCodeGH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GH for simplicity.
	CountryCodeGH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GH
	//CountryCodeGI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GI for simplicity.
	CountryCodeGI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GI
	//CountryCodeGL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GL for simplicity.
	CountryCodeGL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GL
	//CountryCodeGM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GM for simplicity.
	CountryCodeGM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GM
	//CountryCodeGN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GN for simplicity.
	CountryCodeGN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GN
	//CountryCodeGP is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GP for simplicity.
	CountryCodeGP = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GP
	//CountryCodeGQ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GQ for simplicity.
	CountryCodeGQ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GQ
	//CountryCodeGR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GR for simplicity.
	CountryCodeGR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GR
	//CountryCodeGS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GS for simplicity.
	CountryCodeGS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GS
	//CountryCodeGT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GT for simplicity.
	CountryCodeGT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GT
	//CountryCodeGU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GU for simplicity.
	CountryCodeGU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GU
	//CountryCodeGW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GW for simplicity.
	CountryCodeGW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GW
	//CountryCodeGY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GY for simplicity.
	CountryCodeGY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_GY
	//CountryCodeHK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HK for simplicity.
	CountryCodeHK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HK
	//CountryCodeHM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HM for simplicity.
	CountryCodeHM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HM
	//CountryCodeHN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HN for simplicity.
	CountryCodeHN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HN
	//CountryCodeHR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HR for simplicity.
	CountryCodeHR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HR
	//CountryCodeHT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HT for simplicity.
	CountryCodeHT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HT
	//CountryCodeHU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HU for simplicity.
	CountryCodeHU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_HU
	//CountryCodeID is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ID for simplicity.
	CountryCodeID = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ID
	//CountryCodeIE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IE for simplicity.
	CountryCodeIE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IE
	//CountryCodeIL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IL for simplicity.
	CountryCodeIL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IL
	//CountryCodeIM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IM for simplicity.
	CountryCodeIM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IM
	//CountryCodeIN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IN for simplicity.
	CountryCodeIN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IN
	//CountryCodeIO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IO for simplicity.
	CountryCodeIO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IO
	//CountryCodeIQ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IQ for simplicity.
	CountryCodeIQ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IQ
	//CountryCodeIR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IR for simplicity.
	CountryCodeIR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IR
	//CountryCodeIS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IS for simplicity.
	CountryCodeIS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IS
	//CountryCodeIT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IT for simplicity.
	CountryCodeIT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_IT
	//CountryCodeJE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JE for simplicity.
	CountryCodeJE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JE
	//CountryCodeJM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JM for simplicity.
	CountryCodeJM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JM
	//CountryCodeJO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JO for simplicity.
	CountryCodeJO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JO
	//CountryCodeJP is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JP for simplicity.
	CountryCodeJP = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_JP
	//CountryCodeKE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KE for simplicity.
	CountryCodeKE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KE
	//CountryCodeKG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KG for simplicity.
	CountryCodeKG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KG
	//CountryCodeKH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KH for simplicity.
	CountryCodeKH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KH
	//CountryCodeKI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KI for simplicity.
	CountryCodeKI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KI
	//CountryCodeKM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KM for simplicity.
	CountryCodeKM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KM
	//CountryCodeKN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KN for simplicity.
	CountryCodeKN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KN
	//CountryCodeKP is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KP for simplicity.
	CountryCodeKP = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KP
	//CountryCodeKR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KR for simplicity.
	CountryCodeKR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KR
	//CountryCodeKW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KW for simplicity.
	CountryCodeKW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KW
	//CountryCodeKY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KY for simplicity.
	CountryCodeKY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KY
	//CountryCodeKZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KZ for simplicity.
	CountryCodeKZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_KZ
	//CountryCodeLA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LA for simplicity.
	CountryCodeLA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LA
	//CountryCodeLB is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LB for simplicity.
	CountryCodeLB = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LB
	//CountryCodeLC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LC for simplicity.
	CountryCodeLC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LC
	//CountryCodeLI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LI for simplicity.
	CountryCodeLI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LI
	//CountryCodeLK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LK for simplicity.
	CountryCodeLK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LK
	//CountryCodeLR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LR for simplicity.
	CountryCodeLR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LR
	//CountryCodeLS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LS for simplicity.
	CountryCodeLS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LS
	//CountryCodeLT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LT for simplicity.
	CountryCodeLT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LT
	//CountryCodeLU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LU for simplicity.
	CountryCodeLU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LU
	//CountryCodeLV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LV for simplicity.
	CountryCodeLV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LV
	//CountryCodeLY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LY for simplicity.
	CountryCodeLY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_LY
	//CountryCodeMA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MA for simplicity.
	CountryCodeMA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MA
	//CountryCodeMC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MC for simplicity.
	CountryCodeMC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MC
	//CountryCodeMD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MD for simplicity.
	CountryCodeMD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MD
	//CountryCodeME is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ME for simplicity.
	CountryCodeME = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ME
	//CountryCodeMF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MF for simplicity.
	CountryCodeMF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MF
	//CountryCodeMG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MG for simplicity.
	CountryCodeMG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MG
	//CountryCodeMH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MH for simplicity.
	CountryCodeMH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MH
	//CountryCodeMK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MK for simplicity.
	CountryCodeMK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MK
	//CountryCodeML is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ML for simplicity.
	CountryCodeML = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ML
	//CountryCodeMM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MM for simplicity.
	CountryCodeMM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MM
	//CountryCodeMN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MN for simplicity.
	CountryCodeMN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MN
	//CountryCodeMO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MO for simplicity.
	CountryCodeMO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MO
	//CountryCodeMP is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MP for simplicity.
	CountryCodeMP = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MP
	//CountryCodeMQ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MQ for simplicity.
	CountryCodeMQ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MQ
	//CountryCodeMR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MR for simplicity.
	CountryCodeMR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MR
	//CountryCodeMS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MS for simplicity.
	CountryCodeMS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MS
	//CountryCodeMT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MT for simplicity.
	CountryCodeMT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MT
	//CountryCodeMU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MU for simplicity.
	CountryCodeMU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MU
	//CountryCodeMV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MV for simplicity.
	CountryCodeMV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MV
	//CountryCodeMW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MW for simplicity.
	CountryCodeMW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MW
	//CountryCodeMX is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MX for simplicity.
	CountryCodeMX = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MX
	//CountryCodeMY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MY for simplicity.
	CountryCodeMY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MY
	//CountryCodeMZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MZ for simplicity.
	CountryCodeMZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_MZ
	//CountryCodeNA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NA for simplicity.
	CountryCodeNA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NA
	//CountryCodeNC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NC for simplicity.
	CountryCodeNC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NC
	//CountryCodeNE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NE for simplicity.
	CountryCodeNE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NE
	//CountryCodeNF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NF for simplicity.
	CountryCodeNF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NF
	//CountryCodeNG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NG for simplicity.
	CountryCodeNG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NG
	//CountryCodeNI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NI for simplicity.
	CountryCodeNI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NI
	//CountryCodeNL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NL for simplicity.
	CountryCodeNL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NL
	//CountryCodeNO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NO for simplicity.
	CountryCodeNO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NO
	//CountryCodeNP is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NP for simplicity.
	CountryCodeNP = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NP
	//CountryCodeNR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NR for simplicity.
	CountryCodeNR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NR
	//CountryCodeNU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NU for simplicity.
	CountryCodeNU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NU
	//CountryCodeNZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NZ for simplicity.
	CountryCodeNZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NZ
	//CountryCodeOM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_OM for simplicity.
	CountryCodeOM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_OM
	//CountryCodePA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PA for simplicity.
	CountryCodePA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PA
	//CountryCodePE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PE for simplicity.
	CountryCodePE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PE
	//CountryCodePF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PF for simplicity.
	CountryCodePF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PF
	//CountryCodePG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PG for simplicity.
	CountryCodePG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PG
	//CountryCodePH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PH for simplicity.
	CountryCodePH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PH
	//CountryCodePK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PK for simplicity.
	CountryCodePK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PK
	//CountryCodePL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PL for simplicity.
	CountryCodePL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PL
	//CountryCodePM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PM for simplicity.
	CountryCodePM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PM
	//CountryCodePN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PN for simplicity.
	CountryCodePN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PN
	//CountryCodePR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PR for simplicity.
	CountryCodePR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PR
	//CountryCodePS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PS for simplicity.
	CountryCodePS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PS
	//CountryCodePT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PT for simplicity.
	CountryCodePT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PT
	//CountryCodePW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PW for simplicity.
	CountryCodePW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PW
	//CountryCodePY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PY for simplicity.
	CountryCodePY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_PY
	//CountryCodeQA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_QA for simplicity.
	CountryCodeQA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_QA
	//CountryCodeRE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RE for simplicity.
	CountryCodeRE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RE
	//CountryCodeRO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RO for simplicity.
	CountryCodeRO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RO
	//CountryCodeRS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RS for simplicity.
	CountryCodeRS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RS
	//CountryCodeRU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RU for simplicity.
	CountryCodeRU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RU
	//CountryCodeRW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RW for simplicity.
	CountryCodeRW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_RW
	//CountryCodeSA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SA for simplicity.
	CountryCodeSA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SA
	//CountryCodeSB is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SB for simplicity.
	CountryCodeSB = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SB
	//CountryCodeSC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SC for simplicity.
	CountryCodeSC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SC
	//CountryCodeSD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SD for simplicity.
	CountryCodeSD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SD
	//CountryCodeSE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SE for simplicity.
	CountryCodeSE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SE
	//CountryCodeSG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SG for simplicity.
	CountryCodeSG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SG
	//CountryCodeSH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SH for simplicity.
	CountryCodeSH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SH
	//CountryCodeSI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SI for simplicity.
	CountryCodeSI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SI
	//CountryCodeSJ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SJ for simplicity.
	CountryCodeSJ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SJ
	//CountryCodeSK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SK for simplicity.
	CountryCodeSK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SK
	//CountryCodeSL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SL for simplicity.
	CountryCodeSL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SL
	//CountryCodeSM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SM for simplicity.
	CountryCodeSM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SM
	//CountryCodeSN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SN for simplicity.
	CountryCodeSN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SN
	//CountryCodeSO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SO for simplicity.
	CountryCodeSO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SO
	//CountryCodeSR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SR for simplicity.
	CountryCodeSR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SR
	//CountryCodeSS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SS for simplicity.
	CountryCodeSS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SS
	//CountryCodeST is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ST for simplicity.
	CountryCodeST = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ST
	//CountryCodeSV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SV for simplicity.
	CountryCodeSV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SV
	//CountryCodeSX is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SX for simplicity.
	CountryCodeSX = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SX
	//CountryCodeSY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SY for simplicity.
	CountryCodeSY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SY
	//CountryCodeSZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SZ for simplicity.
	CountryCodeSZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_SZ
	//CountryCodeTC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TC for simplicity.
	CountryCodeTC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TC
	//CountryCodeTD is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TD for simplicity.
	CountryCodeTD = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TD
	//CountryCodeTF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TF for simplicity.
	CountryCodeTF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TF
	//CountryCodeTG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TG for simplicity.
	CountryCodeTG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TG
	//CountryCodeTH is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TH for simplicity.
	CountryCodeTH = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TH
	//CountryCodeTJ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TJ for simplicity.
	CountryCodeTJ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TJ
	//CountryCodeTK is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TK for simplicity.
	CountryCodeTK = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TK
	//CountryCodeTL is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TL for simplicity.
	CountryCodeTL = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TL
	//CountryCodeTM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TM for simplicity.
	CountryCodeTM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TM
	//CountryCodeTN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TN for simplicity.
	CountryCodeTN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TN
	//CountryCodeTO is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TO for simplicity.
	CountryCodeTO = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TO
	//CountryCodeTR is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TR for simplicity.
	CountryCodeTR = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TR
	//CountryCodeTT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TT for simplicity.
	CountryCodeTT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TT
	//CountryCodeTV is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TV for simplicity.
	CountryCodeTV = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TV
	//CountryCodeTW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TW for simplicity.
	CountryCodeTW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TW
	//CountryCodeTZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TZ for simplicity.
	CountryCodeTZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_TZ
	//CountryCodeUA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UA for simplicity.
	CountryCodeUA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UA
	//CountryCodeUG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UG for simplicity.
	CountryCodeUG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UG
	//CountryCodeUM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UM for simplicity.
	CountryCodeUM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UM
	//CountryCodeUS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_US for simplicity.
	CountryCodeUS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_US
	//CountryCodeUY is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UY for simplicity.
	CountryCodeUY = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UY
	//CountryCodeUZ is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UZ for simplicity.
	CountryCodeUZ = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_UZ
	//CountryCodeVA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VA for simplicity.
	CountryCodeVA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VA
	//CountryCodeVC is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VC for simplicity.
	CountryCodeVC = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VC
	//CountryCodeVE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VE for simplicity.
	CountryCodeVE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VE
	//CountryCodeVG is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VG for simplicity.
	CountryCodeVG = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VG
	//CountryCodeVI is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VI for simplicity.
	CountryCodeVI = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VI
	//CountryCodeVN is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VN for simplicity.
	CountryCodeVN = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VN
	//CountryCodeVU is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VU for simplicity.
	CountryCodeVU = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_VU
	//CountryCodeWF is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_WF for simplicity.
	CountryCodeWF = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_WF
	//CountryCodeWS is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_WS for simplicity.
	CountryCodeWS = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_WS
	//CountryCodeYE is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_YE for simplicity.
	CountryCodeYE = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_YE
	//CountryCodeYT is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_YT for simplicity.
	CountryCodeYT = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_YT
	//CountryCodeZA is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZA for simplicity.
	CountryCodeZA = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZA
	//CountryCodeZM is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZM for simplicity.
	CountryCodeZM = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZM
	//CountryCodeZW is a simple holder for CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZW for simplicity.
	CountryCodeZW = CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_ZW

	//CountryCodeABW is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ABW for simplicity.
	CountryCodeABW = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ABW
	//CountryCodeAFG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AFG for simplicity.
	CountryCodeAFG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AFG
	//CountryCodeAGO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AGO for simplicity.
	CountryCodeAGO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AGO
	//CountryCodeAIA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AIA for simplicity.
	CountryCodeAIA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AIA
	//CountryCodeALA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ALA for simplicity.
	CountryCodeALA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ALA
	//CountryCodeALB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ALB for simplicity.
	CountryCodeALB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ALB
	//CountryCodeAND is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AND for simplicity.
	CountryCodeAND = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AND
	//CountryCodeARE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARE for simplicity.
	CountryCodeARE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARE
	//CountryCodeARG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARG for simplicity.
	CountryCodeARG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARG
	//CountryCodeARM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARM for simplicity.
	CountryCodeARM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ARM
	//CountryCodeASM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ASM for simplicity.
	CountryCodeASM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ASM
	//CountryCodeATA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATA for simplicity.
	CountryCodeATA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATA
	//CountryCodeATF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATF for simplicity.
	CountryCodeATF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATF
	//CountryCodeATG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATG for simplicity.
	CountryCodeATG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ATG
	//CountryCodeAUS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AUS for simplicity.
	CountryCodeAUS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AUS
	//CountryCodeAUT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AUT for simplicity.
	CountryCodeAUT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AUT
	//CountryCodeAZE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AZE for simplicity.
	CountryCodeAZE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_AZE
	//CountryCodeBDI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BDI for simplicity.
	CountryCodeBDI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BDI
	//CountryCodeBEL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BEL for simplicity.
	CountryCodeBEL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BEL
	//CountryCodeBEN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BEN for simplicity.
	CountryCodeBEN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BEN
	//CountryCodeBES is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BES for simplicity.
	CountryCodeBES = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BES
	//CountryCodeBFA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BFA for simplicity.
	CountryCodeBFA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BFA
	//CountryCodeBGD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BGD for simplicity.
	CountryCodeBGD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BGD
	//CountryCodeBGR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BGR for simplicity.
	CountryCodeBGR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BGR
	//CountryCodeBHR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BHR for simplicity.
	CountryCodeBHR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BHR
	//CountryCodeBHS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BHS for simplicity.
	CountryCodeBHS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BHS
	//CountryCodeBIH is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BIH for simplicity.
	CountryCodeBIH = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BIH
	//CountryCodeBLM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLM for simplicity.
	CountryCodeBLM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLM
	//CountryCodeBLR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLR for simplicity.
	CountryCodeBLR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLR
	//CountryCodeBLZ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLZ for simplicity.
	CountryCodeBLZ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BLZ
	//CountryCodeBMU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BMU for simplicity.
	CountryCodeBMU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BMU
	//CountryCodeBOL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BOL for simplicity.
	CountryCodeBOL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BOL
	//CountryCodeBRA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRA for simplicity.
	CountryCodeBRA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRA
	//CountryCodeBRB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRB for simplicity.
	CountryCodeBRB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRB
	//CountryCodeBRN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRN for simplicity.
	CountryCodeBRN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BRN
	//CountryCodeBTN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BTN for simplicity.
	CountryCodeBTN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BTN
	//CountryCodeBVT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BVT for simplicity.
	CountryCodeBVT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BVT
	//CountryCodeBWA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BWA for simplicity.
	CountryCodeBWA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_BWA
	//CountryCodeCAF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CAF for simplicity.
	CountryCodeCAF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CAF
	//CountryCodeCAN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CAN for simplicity.
	CountryCodeCAN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CAN
	//CountryCodeCCK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CCK for simplicity.
	CountryCodeCCK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CCK
	//CountryCodeCHE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHE for simplicity.
	CountryCodeCHE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHE
	//CountryCodeCHL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHL for simplicity.
	CountryCodeCHL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHL
	//CountryCodeCHN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHN for simplicity.
	CountryCodeCHN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CHN
	//CountryCodeCIV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CIV for simplicity.
	CountryCodeCIV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CIV
	//CountryCodeCMR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CMR for simplicity.
	CountryCodeCMR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CMR
	//CountryCodeCOD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COD for simplicity.
	CountryCodeCOD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COD
	//CountryCodeCOG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COG for simplicity.
	CountryCodeCOG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COG
	//CountryCodeCOK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COK for simplicity.
	CountryCodeCOK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COK
	//CountryCodeCOL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COL for simplicity.
	CountryCodeCOL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COL
	//CountryCodeCOM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COM for simplicity.
	CountryCodeCOM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_COM
	//CountryCodeCPV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CPV for simplicity.
	CountryCodeCPV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CPV
	//CountryCodeCRI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CRI for simplicity.
	CountryCodeCRI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CRI
	//CountryCodeCUB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CUB for simplicity.
	CountryCodeCUB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CUB
	//CountryCodeCUW is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CUW for simplicity.
	CountryCodeCUW = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CUW
	//CountryCodeCXR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CXR for simplicity.
	CountryCodeCXR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CXR
	//CountryCodeCYM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CYM for simplicity.
	CountryCodeCYM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CYM
	//CountryCodeCYP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CYP for simplicity.
	CountryCodeCYP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CYP
	//CountryCodeCZE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CZE for simplicity.
	CountryCodeCZE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_CZE
	//CountryCodeDEU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DEU for simplicity.
	CountryCodeDEU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DEU
	//CountryCodeDJI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DJI for simplicity.
	CountryCodeDJI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DJI
	//CountryCodeDMA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DMA for simplicity.
	CountryCodeDMA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DMA
	//CountryCodeDNK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DNK for simplicity.
	CountryCodeDNK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DNK
	//CountryCodeDOM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DOM for simplicity.
	CountryCodeDOM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DOM
	//CountryCodeDZA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DZA for simplicity.
	CountryCodeDZA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_DZA
	//CountryCodeECU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ECU for simplicity.
	CountryCodeECU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ECU
	//CountryCodeEGY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_EGY for simplicity.
	CountryCodeEGY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_EGY
	//CountryCodeERI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ERI for simplicity.
	CountryCodeERI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ERI
	//CountryCodeESH is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ESH for simplicity.
	CountryCodeESH = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ESH
	//CountryCodeESP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ESP for simplicity.
	CountryCodeESP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ESP
	//CountryCodeEST is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_EST for simplicity.
	CountryCodeEST = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_EST
	//CountryCodeETH is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ETH for simplicity.
	CountryCodeETH = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ETH
	//CountryCodeFIN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FIN for simplicity.
	CountryCodeFIN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FIN
	//CountryCodeFJI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FJI for simplicity.
	CountryCodeFJI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FJI
	//CountryCodeFLK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FLK for simplicity.
	CountryCodeFLK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FLK
	//CountryCodeFRA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FRA for simplicity.
	CountryCodeFRA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FRA
	//CountryCodeFRO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FRO for simplicity.
	CountryCodeFRO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FRO
	//CountryCodeFSM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FSM for simplicity.
	CountryCodeFSM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_FSM
	//CountryCodeGAB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GAB for simplicity.
	CountryCodeGAB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GAB
	//CountryCodeGBR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GBR for simplicity.
	CountryCodeGBR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GBR
	//CountryCodeGEO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GEO for simplicity.
	CountryCodeGEO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GEO
	//CountryCodeGGY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GGY for simplicity.
	CountryCodeGGY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GGY
	//CountryCodeGHA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GHA for simplicity.
	CountryCodeGHA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GHA
	//CountryCodeGIB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GIB for simplicity.
	CountryCodeGIB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GIB
	//CountryCodeGIN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GIN for simplicity.
	CountryCodeGIN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GIN
	//CountryCodeGLP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GLP for simplicity.
	CountryCodeGLP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GLP
	//CountryCodeGMB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GMB for simplicity.
	CountryCodeGMB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GMB
	//CountryCodeGNB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GNB for simplicity.
	CountryCodeGNB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GNB
	//CountryCodeGNQ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GNQ for simplicity.
	CountryCodeGNQ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GNQ
	//CountryCodeGRC is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRC for simplicity.
	CountryCodeGRC = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRC
	//CountryCodeGRD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRD for simplicity.
	CountryCodeGRD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRD
	//CountryCodeGRL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRL for simplicity.
	CountryCodeGRL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GRL
	//CountryCodeGTM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GTM for simplicity.
	CountryCodeGTM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GTM
	//CountryCodeGUF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUF for simplicity.
	CountryCodeGUF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUF
	//CountryCodeGUM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUM for simplicity.
	CountryCodeGUM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUM
	//CountryCodeGUY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUY for simplicity.
	CountryCodeGUY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_GUY
	//CountryCodeHKG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HKG for simplicity.
	CountryCodeHKG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HKG
	//CountryCodeHMD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HMD for simplicity.
	CountryCodeHMD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HMD
	//CountryCodeHND is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HND for simplicity.
	CountryCodeHND = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HND
	//CountryCodeHRV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HRV for simplicity.
	CountryCodeHRV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HRV
	//CountryCodeHTI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HTI for simplicity.
	CountryCodeHTI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HTI
	//CountryCodeHUN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HUN for simplicity.
	CountryCodeHUN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_HUN
	//CountryCodeIDN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IDN for simplicity.
	CountryCodeIDN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IDN
	//CountryCodeIMN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IMN for simplicity.
	CountryCodeIMN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IMN
	//CountryCodeIND is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IND for simplicity.
	CountryCodeIND = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IND
	//CountryCodeIOT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IOT for simplicity.
	CountryCodeIOT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IOT
	//CountryCodeIRL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRL for simplicity.
	CountryCodeIRL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRL
	//CountryCodeIRN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRN for simplicity.
	CountryCodeIRN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRN
	//CountryCodeIRQ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRQ for simplicity.
	CountryCodeIRQ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_IRQ
	//CountryCodeISL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ISL for simplicity.
	CountryCodeISL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ISL
	//CountryCodeISR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ISR for simplicity.
	CountryCodeISR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ISR
	//CountryCodeITA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ITA for simplicity.
	CountryCodeITA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ITA
	//CountryCodeJAM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JAM for simplicity.
	CountryCodeJAM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JAM
	//CountryCodeJEY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JEY for simplicity.
	CountryCodeJEY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JEY
	//CountryCodeJOR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JOR for simplicity.
	CountryCodeJOR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JOR
	//CountryCodeJPN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JPN for simplicity.
	CountryCodeJPN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_JPN
	//CountryCodeKAZ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KAZ for simplicity.
	CountryCodeKAZ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KAZ
	//CountryCodeKEN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KEN for simplicity.
	CountryCodeKEN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KEN
	//CountryCodeKGZ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KGZ for simplicity.
	CountryCodeKGZ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KGZ
	//CountryCodeKHM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KHM for simplicity.
	CountryCodeKHM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KHM
	//CountryCodeKIR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KIR for simplicity.
	CountryCodeKIR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KIR
	//CountryCodeKNA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KNA for simplicity.
	CountryCodeKNA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KNA
	//CountryCodeKOR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KOR for simplicity.
	CountryCodeKOR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KOR
	//CountryCodeKWT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KWT for simplicity.
	CountryCodeKWT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_KWT
	//CountryCodeLAO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LAO for simplicity.
	CountryCodeLAO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LAO
	//CountryCodeLBN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBN for simplicity.
	CountryCodeLBN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBN
	//CountryCodeLBR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBR for simplicity.
	CountryCodeLBR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBR
	//CountryCodeLBY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBY for simplicity.
	CountryCodeLBY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LBY
	//CountryCodeLCA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LCA for simplicity.
	CountryCodeLCA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LCA
	//CountryCodeLIE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LIE for simplicity.
	CountryCodeLIE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LIE
	//CountryCodeLKA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LKA for simplicity.
	CountryCodeLKA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LKA
	//CountryCodeLSO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LSO for simplicity.
	CountryCodeLSO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LSO
	//CountryCodeLTU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LTU for simplicity.
	CountryCodeLTU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LTU
	//CountryCodeLUX is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LUX for simplicity.
	CountryCodeLUX = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LUX
	//CountryCodeLVA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LVA for simplicity.
	CountryCodeLVA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_LVA
	//CountryCodeMAC is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAC for simplicity.
	CountryCodeMAC = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAC
	//CountryCodeMAF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAF for simplicity.
	CountryCodeMAF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAF
	//CountryCodeMAR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAR for simplicity.
	CountryCodeMAR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MAR
	//CountryCodeMCO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MCO for simplicity.
	CountryCodeMCO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MCO
	//CountryCodeMDA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDA for simplicity.
	CountryCodeMDA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDA
	//CountryCodeMDG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDG for simplicity.
	CountryCodeMDG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDG
	//CountryCodeMDV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDV for simplicity.
	CountryCodeMDV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MDV
	//CountryCodeMEX is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MEX for simplicity.
	CountryCodeMEX = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MEX
	//CountryCodeMHL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MHL for simplicity.
	CountryCodeMHL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MHL
	//CountryCodeMKD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MKD for simplicity.
	CountryCodeMKD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MKD
	//CountryCodeMLI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MLI for simplicity.
	CountryCodeMLI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MLI
	//CountryCodeMLT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MLT for simplicity.
	CountryCodeMLT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MLT
	//CountryCodeMMR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MMR for simplicity.
	CountryCodeMMR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MMR
	//CountryCodeMNE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNE for simplicity.
	CountryCodeMNE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNE
	//CountryCodeMNG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNG for simplicity.
	CountryCodeMNG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNG
	//CountryCodeMNP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNP for simplicity.
	CountryCodeMNP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MNP
	//CountryCodeMOZ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MOZ for simplicity.
	CountryCodeMOZ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MOZ
	//CountryCodeMRT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MRT for simplicity.
	CountryCodeMRT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MRT
	//CountryCodeMSR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MSR for simplicity.
	CountryCodeMSR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MSR
	//CountryCodeMTQ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MTQ for simplicity.
	CountryCodeMTQ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MTQ
	//CountryCodeMUS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MUS for simplicity.
	CountryCodeMUS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MUS
	//CountryCodeMWI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MWI for simplicity.
	CountryCodeMWI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MWI
	//CountryCodeMYS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MYS for simplicity.
	CountryCodeMYS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MYS
	//CountryCodeMYT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MYT for simplicity.
	CountryCodeMYT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_MYT
	//CountryCodeNAM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NAM for simplicity.
	CountryCodeNAM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NAM
	//CountryCodeNCL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NCL for simplicity.
	CountryCodeNCL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NCL
	//CountryCodeNER is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NER for simplicity.
	CountryCodeNER = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NER
	//CountryCodeNFK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NFK for simplicity.
	CountryCodeNFK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NFK
	//CountryCodeNGA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NGA for simplicity.
	CountryCodeNGA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NGA
	//CountryCodeNIC is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NIC for simplicity.
	CountryCodeNIC = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NIC
	//CountryCodeNIU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NIU for simplicity.
	CountryCodeNIU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NIU
	//CountryCodeNLD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NLD for simplicity.
	CountryCodeNLD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NLD
	//CountryCodeNOR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NOR for simplicity.
	CountryCodeNOR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NOR
	//CountryCodeNPL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NPL for simplicity.
	CountryCodeNPL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NPL
	//CountryCodeNRU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NRU for simplicity.
	CountryCodeNRU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NRU
	//CountryCodeNZL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NZL for simplicity.
	CountryCodeNZL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NZL
	//CountryCodeOMN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_OMN for simplicity.
	CountryCodeOMN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_OMN
	//CountryCodePAK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PAK for simplicity.
	CountryCodePAK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PAK
	//CountryCodePAN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PAN for simplicity.
	CountryCodePAN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PAN
	//CountryCodePCN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PCN for simplicity.
	CountryCodePCN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PCN
	//CountryCodePER is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PER for simplicity.
	CountryCodePER = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PER
	//CountryCodePHL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PHL for simplicity.
	CountryCodePHL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PHL
	//CountryCodePLW is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PLW for simplicity.
	CountryCodePLW = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PLW
	//CountryCodePNG is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PNG for simplicity.
	CountryCodePNG = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PNG
	//CountryCodePOL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_POL for simplicity.
	CountryCodePOL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_POL
	//CountryCodePRI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRI for simplicity.
	CountryCodePRI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRI
	//CountryCodePRK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRK for simplicity.
	CountryCodePRK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRK
	//CountryCodePRT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRT for simplicity.
	CountryCodePRT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRT
	//CountryCodePRY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRY for simplicity.
	CountryCodePRY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PRY
	//CountryCodePSE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PSE for simplicity.
	CountryCodePSE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PSE
	//CountryCodePYF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PYF for simplicity.
	CountryCodePYF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_PYF
	//CountryCodeQAT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_QAT for simplicity.
	CountryCodeQAT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_QAT
	//CountryCodeREU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_REU for simplicity.
	CountryCodeREU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_REU
	//CountryCodeROU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ROU for simplicity.
	CountryCodeROU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ROU
	//CountryCodeRUS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_RUS for simplicity.
	CountryCodeRUS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_RUS
	//CountryCodeRWA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_RWA for simplicity.
	CountryCodeRWA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_RWA
	//CountryCodeSAU is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SAU for simplicity.
	CountryCodeSAU = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SAU
	//CountryCodeSDN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SDN for simplicity.
	CountryCodeSDN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SDN
	//CountryCodeSEN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SEN for simplicity.
	CountryCodeSEN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SEN
	//CountryCodeSGP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SGP for simplicity.
	CountryCodeSGP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SGP
	//CountryCodeSGS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SGS for simplicity.
	CountryCodeSGS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SGS
	//CountryCodeSHN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SHN for simplicity.
	CountryCodeSHN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SHN
	//CountryCodeSJM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SJM for simplicity.
	CountryCodeSJM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SJM
	//CountryCodeSLB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLB for simplicity.
	CountryCodeSLB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLB
	//CountryCodeSLE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLE for simplicity.
	CountryCodeSLE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLE
	//CountryCodeSLV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLV for simplicity.
	CountryCodeSLV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SLV
	//CountryCodeSMR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SMR for simplicity.
	CountryCodeSMR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SMR
	//CountryCodeSOM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SOM for simplicity.
	CountryCodeSOM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SOM
	//CountryCodeSPM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SPM for simplicity.
	CountryCodeSPM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SPM
	//CountryCodeSRB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SRB for simplicity.
	CountryCodeSRB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SRB
	//CountryCodeSSD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SSD for simplicity.
	CountryCodeSSD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SSD
	//CountryCodeSTP is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_STP for simplicity.
	CountryCodeSTP = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_STP
	//CountryCodeSUR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SUR for simplicity.
	CountryCodeSUR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SUR
	//CountryCodeSVK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SVK for simplicity.
	CountryCodeSVK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SVK
	//CountryCodeSVN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SVN for simplicity.
	CountryCodeSVN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SVN
	//CountryCodeSWE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SWE for simplicity.
	CountryCodeSWE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SWE
	//CountryCodeSWZ is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SWZ for simplicity.
	CountryCodeSWZ = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SWZ
	//CountryCodeSXM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SXM for simplicity.
	CountryCodeSXM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SXM
	//CountryCodeSYC is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SYC for simplicity.
	CountryCodeSYC = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SYC
	//CountryCodeSYR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SYR for simplicity.
	CountryCodeSYR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_SYR
	//CountryCodeTCA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TCA for simplicity.
	CountryCodeTCA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TCA
	//CountryCodeTCD is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TCD for simplicity.
	CountryCodeTCD = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TCD
	//CountryCodeTGO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TGO for simplicity.
	CountryCodeTGO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TGO
	//CountryCodeTHA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_THA for simplicity.
	CountryCodeTHA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_THA
	//CountryCodeTJK is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TJK for simplicity.
	CountryCodeTJK = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TJK
	//CountryCodeTKL is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TKL for simplicity.
	CountryCodeTKL = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TKL
	//CountryCodeTKM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TKM for simplicity.
	CountryCodeTKM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TKM
	//CountryCodeTLS is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TLS for simplicity.
	CountryCodeTLS = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TLS
	//CountryCodeTON is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TON for simplicity.
	CountryCodeTON = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TON
	//CountryCodeTTO is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TTO for simplicity.
	CountryCodeTTO = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TTO
	//CountryCodeTUN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUN for simplicity.
	CountryCodeTUN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUN
	//CountryCodeTUR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUR for simplicity.
	CountryCodeTUR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUR
	//CountryCodeTUV is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUV for simplicity.
	CountryCodeTUV = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TUV
	//CountryCodeTWN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TWN for simplicity.
	CountryCodeTWN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TWN
	//CountryCodeTZA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TZA for simplicity.
	CountryCodeTZA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_TZA
	//CountryCodeUGA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UGA for simplicity.
	CountryCodeUGA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UGA
	//CountryCodeUKR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UKR for simplicity.
	CountryCodeUKR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UKR
	//CountryCodeUMI is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UMI for simplicity.
	CountryCodeUMI = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UMI
	//CountryCodeURY is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_URY for simplicity.
	CountryCodeURY = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_URY
	//CountryCodeUSA is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_USA for simplicity.
	CountryCodeUSA = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_USA
	//CountryCodeUZB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UZB for simplicity.
	CountryCodeUZB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_UZB
	//CountryCodeVAT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VAT for simplicity.
	CountryCodeVAT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VAT
	//CountryCodeVCT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VCT for simplicity.
	CountryCodeVCT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VCT
	//CountryCodeVEN is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VEN for simplicity.
	CountryCodeVEN = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VEN
	//CountryCodeVGB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VGB for simplicity.
	CountryCodeVGB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VGB
	//CountryCodeVIR is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VIR for simplicity.
	CountryCodeVIR = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VIR
	//CountryCodeVNM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VNM for simplicity.
	CountryCodeVNM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VNM
	//CountryCodeVUT is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VUT for simplicity.
	CountryCodeVUT = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_VUT
	//CountryCodeWLF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_WLF for simplicity.
	CountryCodeWLF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_WLF
	//CountryCodeWSM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_WSM for simplicity.
	CountryCodeWSM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_WSM
	//CountryCodeYEM is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_YEM for simplicity.
	CountryCodeYEM = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_YEM
	//CountryCodeZAF is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZAF for simplicity.
	CountryCodeZAF = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZAF
	//CountryCodeZMB is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZMB for simplicity.
	CountryCodeZMB = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZMB
	//CountryCodeZWE is a simple holder for CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZWE for simplicity.
	CountryCodeZWE = CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_ZWE

	// CountryAlpha2CodeToCountry is a map from CountryAlpha2Code to Country.
	CountryAlpha2CodeToCountry = map[CountryAlpha2Code]*Country{
		CountryCodeAD: country1,
		CountryCodeAE: country2,
		CountryCodeAF: country3,
		CountryCodeAG: country4,
		CountryCodeAI: country5,
		CountryCodeAL: country6,
		CountryCodeAM: country7,
		CountryCodeAO: country8,
		CountryCodeAQ: country9,
		CountryCodeAR: country10,
		CountryCodeAS: country11,
		CountryCodeAT: country12,
		CountryCodeAU: country13,
		CountryCodeAW: country14,
		CountryCodeAX: country15,
		CountryCodeAZ: country16,
		CountryCodeBA: country17,
		CountryCodeBB: country18,
		CountryCodeBD: country19,
		CountryCodeBE: country20,
		CountryCodeBF: country21,
		CountryCodeBG: country22,
		CountryCodeBH: country23,
		CountryCodeBI: country24,
		CountryCodeBJ: country25,
		CountryCodeBL: country26,
		CountryCodeBM: country27,
		CountryCodeBN: country28,
		CountryCodeBO: country29,
		CountryCodeBQ: country30,
		CountryCodeBR: country31,
		CountryCodeBS: country32,
		CountryCodeBT: country33,
		CountryCodeBV: country34,
		CountryCodeBW: country35,
		CountryCodeBY: country36,
		CountryCodeBZ: country37,
		CountryCodeCA: country38,
		CountryCodeCC: country39,
		CountryCodeCD: country40,
		CountryCodeCF: country41,
		CountryCodeCG: country42,
		CountryCodeCH: country43,
		CountryCodeCI: country44,
		CountryCodeCK: country45,
		CountryCodeCL: country46,
		CountryCodeCM: country47,
		CountryCodeCN: country48,
		CountryCodeCO: country49,
		CountryCodeCR: country50,
		CountryCodeCU: country51,
		CountryCodeCV: country52,
		CountryCodeCW: country53,
		CountryCodeCX: country54,
		CountryCodeCY: country55,
		CountryCodeCZ: country56,
		CountryCodeDE: country57,
		CountryCodeDJ: country58,
		CountryCodeDK: country59,
		CountryCodeDM: country60,
		CountryCodeDO: country61,
		CountryCodeDZ: country62,
		CountryCodeEC: country63,
		CountryCodeEE: country64,
		CountryCodeEG: country65,
		CountryCodeEH: country66,
		CountryCodeER: country67,
		CountryCodeES: country68,
		CountryCodeET: country69,
		CountryCodeFI: country70,
		CountryCodeFJ: country71,
		CountryCodeFK: country72,
		CountryCodeFM: country73,
		CountryCodeFO: country74,
		CountryCodeFR: country75,
		CountryCodeGA: country76,
		CountryCodeGB: country77,
		CountryCodeGD: country78,
		CountryCodeGE: country79,
		CountryCodeGF: country80,
		CountryCodeGG: country81,
		CountryCodeGH: country82,
		CountryCodeGI: country83,
		CountryCodeGL: country84,
		CountryCodeGM: country85,
		CountryCodeGN: country86,
		CountryCodeGP: country87,
		CountryCodeGQ: country88,
		CountryCodeGR: country89,
		CountryCodeGS: country90,
		CountryCodeGT: country91,
		CountryCodeGU: country92,
		CountryCodeGW: country93,
		CountryCodeGY: country94,
		CountryCodeHK: country95,
		CountryCodeHM: country96,
		CountryCodeHN: country97,
		CountryCodeHR: country98,
		CountryCodeHT: country99,
		CountryCodeHU: country100,
		CountryCodeID: country101,
		CountryCodeIE: country102,
		CountryCodeIL: country103,
		CountryCodeIM: country104,
		CountryCodeIN: country105,
		CountryCodeIO: country106,
		CountryCodeIQ: country107,
		CountryCodeIR: country108,
		CountryCodeIS: country109,
		CountryCodeIT: country110,
		CountryCodeJE: country111,
		CountryCodeJM: country112,
		CountryCodeJO: country113,
		CountryCodeJP: country114,
		CountryCodeKE: country115,
		CountryCodeKG: country116,
		CountryCodeKH: country117,
		CountryCodeKI: country118,
		CountryCodeKM: country119,
		CountryCodeKN: country120,
		CountryCodeKP: country121,
		CountryCodeKR: country122,
		CountryCodeKW: country123,
		CountryCodeKY: country124,
		CountryCodeKZ: country125,
		CountryCodeLA: country126,
		CountryCodeLB: country127,
		CountryCodeLC: country128,
		CountryCodeLI: country129,
		CountryCodeLK: country130,
		CountryCodeLR: country131,
		CountryCodeLS: country132,
		CountryCodeLT: country133,
		CountryCodeLU: country134,
		CountryCodeLV: country135,
		CountryCodeLY: country136,
		CountryCodeMA: country137,
		CountryCodeMC: country138,
		CountryCodeMD: country139,
		CountryCodeME: country140,
		CountryCodeMF: country141,
		CountryCodeMG: country142,
		CountryCodeMH: country143,
		CountryCodeMK: country144,
		CountryCodeML: country145,
		CountryCodeMM: country146,
		CountryCodeMN: country147,
		CountryCodeMO: country148,
		CountryCodeMP: country149,
		CountryCodeMQ: country150,
		CountryCodeMR: country151,
		CountryCodeMS: country152,
		CountryCodeMT: country153,
		CountryCodeMU: country154,
		CountryCodeMV: country155,
		CountryCodeMW: country156,
		CountryCodeMX: country157,
		CountryCodeMY: country158,
		CountryCodeMZ: country159,
		CountryCodeNA: country160,
		CountryCodeNC: country161,
		CountryCodeNE: country162,
		CountryCodeNF: country163,
		CountryCodeNG: country164,
		CountryCodeNI: country165,
		CountryCodeNL: country166,
		CountryCodeNO: country167,
		CountryCodeNP: country168,
		CountryCodeNR: country169,
		CountryCodeNU: country170,
		CountryCodeNZ: country171,
		CountryCodeOM: country172,
		CountryCodePA: country173,
		CountryCodePE: country174,
		CountryCodePF: country175,
		CountryCodePG: country176,
		CountryCodePH: country177,
		CountryCodePK: country178,
		CountryCodePL: country179,
		CountryCodePM: country180,
		CountryCodePN: country181,
		CountryCodePR: country182,
		CountryCodePS: country183,
		CountryCodePT: country184,
		CountryCodePW: country185,
		CountryCodePY: country186,
		CountryCodeQA: country187,
		CountryCodeRE: country188,
		CountryCodeRO: country189,
		CountryCodeRS: country190,
		CountryCodeRU: country191,
		CountryCodeRW: country192,
		CountryCodeSA: country193,
		CountryCodeSB: country194,
		CountryCodeSC: country195,
		CountryCodeSD: country196,
		CountryCodeSE: country197,
		CountryCodeSG: country198,
		CountryCodeSH: country199,
		CountryCodeSI: country200,
		CountryCodeSJ: country201,
		CountryCodeSK: country202,
		CountryCodeSL: country203,
		CountryCodeSM: country204,
		CountryCodeSN: country205,
		CountryCodeSO: country206,
		CountryCodeSR: country207,
		CountryCodeSS: country208,
		CountryCodeST: country209,
		CountryCodeSV: country210,
		CountryCodeSX: country211,
		CountryCodeSY: country212,
		CountryCodeSZ: country213,
		CountryCodeTC: country214,
		CountryCodeTD: country215,
		CountryCodeTF: country216,
		CountryCodeTG: country217,
		CountryCodeTH: country218,
		CountryCodeTJ: country219,
		CountryCodeTK: country220,
		CountryCodeTL: country221,
		CountryCodeTM: country222,
		CountryCodeTN: country223,
		CountryCodeTO: country224,
		CountryCodeTR: country225,
		CountryCodeTT: country226,
		CountryCodeTV: country227,
		CountryCodeTW: country228,
		CountryCodeTZ: country229,
		CountryCodeUA: country230,
		CountryCodeUG: country231,
		CountryCodeUM: country232,
		CountryCodeUS: country233,
		CountryCodeUY: country234,
		CountryCodeUZ: country235,
		CountryCodeVA: country236,
		CountryCodeVC: country237,
		CountryCodeVE: country238,
		CountryCodeVG: country239,
		CountryCodeVI: country240,
		CountryCodeVN: country241,
		CountryCodeVU: country242,
		CountryCodeWF: country243,
		CountryCodeWS: country244,
		CountryCodeYE: country245,
		CountryCodeYT: country246,
		CountryCodeZA: country247,
		CountryCodeZM: country248,
		CountryCodeZW: country249,
	}

	// CountryAlpha3CodeToCountry is a map from CountryAlpha3Code to Country.
	CountryAlpha3CodeToCountry = map[CountryAlpha3Code]*Country{
		CountryCodeABW: country1,
		CountryCodeAFG: country2,
		CountryCodeAGO: country3,
		CountryCodeAIA: country4,
		CountryCodeALA: country5,
		CountryCodeALB: country6,
		CountryCodeAND: country7,
		CountryCodeARE: country8,
		CountryCodeARG: country9,
		CountryCodeARM: country10,
		CountryCodeASM: country11,
		CountryCodeATA: country12,
		CountryCodeATF: country13,
		CountryCodeATG: country14,
		CountryCodeAUS: country15,
		CountryCodeAUT: country16,
		CountryCodeAZE: country17,
		CountryCodeBDI: country18,
		CountryCodeBEL: country19,
		CountryCodeBEN: country20,
		CountryCodeBES: country21,
		CountryCodeBFA: country22,
		CountryCodeBGD: country23,
		CountryCodeBGR: country24,
		CountryCodeBHR: country25,
		CountryCodeBHS: country26,
		CountryCodeBIH: country27,
		CountryCodeBLM: country28,
		CountryCodeBLR: country29,
		CountryCodeBLZ: country30,
		CountryCodeBMU: country31,
		CountryCodeBOL: country32,
		CountryCodeBRA: country33,
		CountryCodeBRB: country34,
		CountryCodeBRN: country35,
		CountryCodeBTN: country36,
		CountryCodeBVT: country37,
		CountryCodeBWA: country38,
		CountryCodeCAF: country39,
		CountryCodeCAN: country40,
		CountryCodeCCK: country41,
		CountryCodeCHE: country42,
		CountryCodeCHL: country43,
		CountryCodeCHN: country44,
		CountryCodeCIV: country45,
		CountryCodeCMR: country46,
		CountryCodeCOD: country47,
		CountryCodeCOG: country48,
		CountryCodeCOK: country49,
		CountryCodeCOL: country50,
		CountryCodeCOM: country51,
		CountryCodeCPV: country52,
		CountryCodeCRI: country53,
		CountryCodeCUB: country54,
		CountryCodeCUW: country55,
		CountryCodeCXR: country56,
		CountryCodeCYM: country57,
		CountryCodeCYP: country58,
		CountryCodeCZE: country59,
		CountryCodeDEU: country60,
		CountryCodeDJI: country61,
		CountryCodeDMA: country62,
		CountryCodeDNK: country63,
		CountryCodeDOM: country64,
		CountryCodeDZA: country65,
		CountryCodeECU: country66,
		CountryCodeEGY: country67,
		CountryCodeERI: country68,
		CountryCodeESH: country69,
		CountryCodeESP: country70,
		CountryCodeEST: country71,
		CountryCodeETH: country72,
		CountryCodeFIN: country73,
		CountryCodeFJI: country74,
		CountryCodeFLK: country75,
		CountryCodeFRA: country76,
		CountryCodeFRO: country77,
		CountryCodeFSM: country78,
		CountryCodeGAB: country79,
		CountryCodeGBR: country80,
		CountryCodeGEO: country81,
		CountryCodeGGY: country82,
		CountryCodeGHA: country83,
		CountryCodeGIB: country84,
		CountryCodeGIN: country85,
		CountryCodeGLP: country86,
		CountryCodeGMB: country87,
		CountryCodeGNB: country88,
		CountryCodeGNQ: country89,
		CountryCodeGRC: country90,
		CountryCodeGRD: country91,
		CountryCodeGRL: country92,
		CountryCodeGTM: country93,
		CountryCodeGUF: country94,
		CountryCodeGUM: country95,
		CountryCodeGUY: country96,
		CountryCodeHKG: country97,
		CountryCodeHMD: country98,
		CountryCodeHND: country99,
		CountryCodeHRV: country100,
		CountryCodeHTI: country101,
		CountryCodeHUN: country102,
		CountryCodeIDN: country103,
		CountryCodeIMN: country104,
		CountryCodeIND: country105,
		CountryCodeIOT: country106,
		CountryCodeIRL: country107,
		CountryCodeIRN: country108,
		CountryCodeIRQ: country109,
		CountryCodeISL: country110,
		CountryCodeISR: country111,
		CountryCodeITA: country112,
		CountryCodeJAM: country113,
		CountryCodeJEY: country114,
		CountryCodeJOR: country115,
		CountryCodeJPN: country116,
		CountryCodeKAZ: country117,
		CountryCodeKEN: country118,
		CountryCodeKGZ: country119,
		CountryCodeKHM: country120,
		CountryCodeKIR: country121,
		CountryCodeKNA: country122,
		CountryCodeKOR: country123,
		CountryCodeKWT: country124,
		CountryCodeLAO: country125,
		CountryCodeLBN: country126,
		CountryCodeLBR: country127,
		CountryCodeLBY: country128,
		CountryCodeLCA: country129,
		CountryCodeLIE: country130,
		CountryCodeLKA: country131,
		CountryCodeLSO: country132,
		CountryCodeLTU: country133,
		CountryCodeLUX: country134,
		CountryCodeLVA: country135,
		CountryCodeMAC: country136,
		CountryCodeMAF: country137,
		CountryCodeMAR: country138,
		CountryCodeMCO: country139,
		CountryCodeMDA: country140,
		CountryCodeMDG: country141,
		CountryCodeMDV: country142,
		CountryCodeMEX: country143,
		CountryCodeMHL: country144,
		CountryCodeMKD: country145,
		CountryCodeMLI: country146,
		CountryCodeMLT: country147,
		CountryCodeMMR: country148,
		CountryCodeMNE: country149,
		CountryCodeMNG: country150,
		CountryCodeMNP: country151,
		CountryCodeMOZ: country152,
		CountryCodeMRT: country153,
		CountryCodeMSR: country154,
		CountryCodeMTQ: country155,
		CountryCodeMUS: country156,
		CountryCodeMWI: country157,
		CountryCodeMYS: country158,
		CountryCodeMYT: country159,
		CountryCodeNAM: country160,
		CountryCodeNCL: country161,
		CountryCodeNER: country162,
		CountryCodeNFK: country163,
		CountryCodeNGA: country164,
		CountryCodeNIC: country165,
		CountryCodeNIU: country166,
		CountryCodeNLD: country167,
		CountryCodeNOR: country168,
		CountryCodeNPL: country169,
		CountryCodeNRU: country170,
		CountryCodeNZL: country171,
		CountryCodeOMN: country172,
		CountryCodePAK: country173,
		CountryCodePAN: country174,
		CountryCodePCN: country175,
		CountryCodePER: country176,
		CountryCodePHL: country177,
		CountryCodePLW: country178,
		CountryCodePNG: country179,
		CountryCodePOL: country180,
		CountryCodePRI: country181,
		CountryCodePRK: country182,
		CountryCodePRT: country183,
		CountryCodePRY: country184,
		CountryCodePSE: country185,
		CountryCodePYF: country186,
		CountryCodeQAT: country187,
		CountryCodeREU: country188,
		CountryCodeROU: country189,
		CountryCodeRUS: country190,
		CountryCodeRWA: country191,
		CountryCodeSAU: country192,
		CountryCodeSDN: country193,
		CountryCodeSEN: country194,
		CountryCodeSGP: country195,
		CountryCodeSGS: country196,
		CountryCodeSHN: country197,
		CountryCodeSJM: country198,
		CountryCodeSLB: country199,
		CountryCodeSLE: country200,
		CountryCodeSLV: country201,
		CountryCodeSMR: country202,
		CountryCodeSOM: country203,
		CountryCodeSPM: country204,
		CountryCodeSRB: country205,
		CountryCodeSSD: country206,
		CountryCodeSTP: country207,
		CountryCodeSUR: country208,
		CountryCodeSVK: country209,
		CountryCodeSVN: country210,
		CountryCodeSWE: country211,
		CountryCodeSWZ: country212,
		CountryCodeSXM: country213,
		CountryCodeSYC: country214,
		CountryCodeSYR: country215,
		CountryCodeTCA: country216,
		CountryCodeTCD: country217,
		CountryCodeTGO: country218,
		CountryCodeTHA: country219,
		CountryCodeTJK: country220,
		CountryCodeTKL: country221,
		CountryCodeTKM: country222,
		CountryCodeTLS: country223,
		CountryCodeTON: country224,
		CountryCodeTTO: country225,
		CountryCodeTUN: country226,
		CountryCodeTUR: country227,
		CountryCodeTUV: country228,
		CountryCodeTWN: country229,
		CountryCodeTZA: country230,
		CountryCodeUGA: country231,
		CountryCodeUKR: country232,
		CountryCodeUMI: country233,
		CountryCodeURY: country234,
		CountryCodeUSA: country235,
		CountryCodeUZB: country236,
		CountryCodeVAT: country237,
		CountryCodeVCT: country238,
		CountryCodeVEN: country239,
		CountryCodeVGB: country240,
		CountryCodeVIR: country241,
		CountryCodeVNM: country242,
		CountryCodeVUT: country243,
		CountryCodeWLF: country244,
		CountryCodeWSM: country245,
		CountryCodeYEM: country246,
		CountryCodeZAF: country247,
		CountryCodeZMB: country248,
		CountryCodeZWE: country249,
	}

	// CountryAlpha2CodeToSimpleString is a map from CountryAlpha2Code to simple string.
	CountryAlpha2CodeToSimpleString = map[CountryAlpha2Code]string{
		CountryCodeAD: "AD",
		CountryCodeAE: "AE",
		CountryCodeAF: "AF",
		CountryCodeAG: "AG",
		CountryCodeAI: "AI",
		CountryCodeAL: "AL",
		CountryCodeAM: "AM",
		CountryCodeAO: "AO",
		CountryCodeAQ: "AQ",
		CountryCodeAR: "AR",
		CountryCodeAS: "AS",
		CountryCodeAT: "AT",
		CountryCodeAU: "AU",
		CountryCodeAW: "AW",
		CountryCodeAX: "AX",
		CountryCodeAZ: "AZ",
		CountryCodeBA: "BA",
		CountryCodeBB: "BB",
		CountryCodeBD: "BD",
		CountryCodeBE: "BE",
		CountryCodeBF: "BF",
		CountryCodeBG: "BG",
		CountryCodeBH: "BH",
		CountryCodeBI: "BI",
		CountryCodeBJ: "BJ",
		CountryCodeBL: "BL",
		CountryCodeBM: "BM",
		CountryCodeBN: "BN",
		CountryCodeBO: "BO",
		CountryCodeBQ: "BQ",
		CountryCodeBR: "BR",
		CountryCodeBS: "BS",
		CountryCodeBT: "BT",
		CountryCodeBV: "BV",
		CountryCodeBW: "BW",
		CountryCodeBY: "BY",
		CountryCodeBZ: "BZ",
		CountryCodeCA: "CA",
		CountryCodeCC: "CC",
		CountryCodeCD: "CD",
		CountryCodeCF: "CF",
		CountryCodeCG: "CG",
		CountryCodeCH: "CH",
		CountryCodeCI: "CI",
		CountryCodeCK: "CK",
		CountryCodeCL: "CL",
		CountryCodeCM: "CM",
		CountryCodeCN: "CN",
		CountryCodeCO: "CO",
		CountryCodeCR: "CR",
		CountryCodeCU: "CU",
		CountryCodeCV: "CV",
		CountryCodeCW: "CW",
		CountryCodeCX: "CX",
		CountryCodeCY: "CY",
		CountryCodeCZ: "CZ",
		CountryCodeDE: "DE",
		CountryCodeDJ: "DJ",
		CountryCodeDK: "DK",
		CountryCodeDM: "DM",
		CountryCodeDO: "DO",
		CountryCodeDZ: "DZ",
		CountryCodeEC: "EC",
		CountryCodeEE: "EE",
		CountryCodeEG: "EG",
		CountryCodeEH: "EH",
		CountryCodeER: "ER",
		CountryCodeES: "ES",
		CountryCodeET: "ET",
		CountryCodeFI: "FI",
		CountryCodeFJ: "FJ",
		CountryCodeFK: "FK",
		CountryCodeFM: "FM",
		CountryCodeFO: "FO",
		CountryCodeFR: "FR",
		CountryCodeGA: "GA",
		CountryCodeGB: "GB",
		CountryCodeGD: "GD",
		CountryCodeGE: "GE",
		CountryCodeGF: "GF",
		CountryCodeGG: "GG",
		CountryCodeGH: "GH",
		CountryCodeGI: "GI",
		CountryCodeGL: "GL",
		CountryCodeGM: "GM",
		CountryCodeGN: "GN",
		CountryCodeGP: "GP",
		CountryCodeGQ: "GQ",
		CountryCodeGR: "GR",
		CountryCodeGS: "GS",
		CountryCodeGT: "GT",
		CountryCodeGU: "GU",
		CountryCodeGW: "GW",
		CountryCodeGY: "GY",
		CountryCodeHK: "HK",
		CountryCodeHM: "HM",
		CountryCodeHN: "HN",
		CountryCodeHR: "HR",
		CountryCodeHT: "HT",
		CountryCodeHU: "HU",
		CountryCodeID: "ID",
		CountryCodeIE: "IE",
		CountryCodeIL: "IL",
		CountryCodeIM: "IM",
		CountryCodeIN: "IN",
		CountryCodeIO: "IO",
		CountryCodeIQ: "IQ",
		CountryCodeIR: "IR",
		CountryCodeIS: "IS",
		CountryCodeIT: "IT",
		CountryCodeJE: "JE",
		CountryCodeJM: "JM",
		CountryCodeJO: "JO",
		CountryCodeJP: "JP",
		CountryCodeKE: "KE",
		CountryCodeKG: "KG",
		CountryCodeKH: "KH",
		CountryCodeKI: "KI",
		CountryCodeKM: "KM",
		CountryCodeKN: "KN",
		CountryCodeKP: "KP",
		CountryCodeKR: "KR",
		CountryCodeKW: "KW",
		CountryCodeKY: "KY",
		CountryCodeKZ: "KZ",
		CountryCodeLA: "LA",
		CountryCodeLB: "LB",
		CountryCodeLC: "LC",
		CountryCodeLI: "LI",
		CountryCodeLK: "LK",
		CountryCodeLR: "LR",
		CountryCodeLS: "LS",
		CountryCodeLT: "LT",
		CountryCodeLU: "LU",
		CountryCodeLV: "LV",
		CountryCodeLY: "LY",
		CountryCodeMA: "MA",
		CountryCodeMC: "MC",
		CountryCodeMD: "MD",
		CountryCodeME: "ME",
		CountryCodeMF: "MF",
		CountryCodeMG: "MG",
		CountryCodeMH: "MH",
		CountryCodeMK: "MK",
		CountryCodeML: "ML",
		CountryCodeMM: "MM",
		CountryCodeMN: "MN",
		CountryCodeMO: "MO",
		CountryCodeMP: "MP",
		CountryCodeMQ: "MQ",
		CountryCodeMR: "MR",
		CountryCodeMS: "MS",
		CountryCodeMT: "MT",
		CountryCodeMU: "MU",
		CountryCodeMV: "MV",
		CountryCodeMW: "MW",
		CountryCodeMX: "MX",
		CountryCodeMY: "MY",
		CountryCodeMZ: "MZ",
		CountryCodeNA: "NA",
		CountryCodeNC: "NC",
		CountryCodeNE: "NE",
		CountryCodeNF: "NF",
		CountryCodeNG: "NG",
		CountryCodeNI: "NI",
		CountryCodeNL: "NL",
		CountryCodeNO: "NO",
		CountryCodeNP: "NP",
		CountryCodeNR: "NR",
		CountryCodeNU: "NU",
		CountryCodeNZ: "NZ",
		CountryCodeOM: "OM",
		CountryCodePA: "PA",
		CountryCodePE: "PE",
		CountryCodePF: "PF",
		CountryCodePG: "PG",
		CountryCodePH: "PH",
		CountryCodePK: "PK",
		CountryCodePL: "PL",
		CountryCodePM: "PM",
		CountryCodePN: "PN",
		CountryCodePR: "PR",
		CountryCodePS: "PS",
		CountryCodePT: "PT",
		CountryCodePW: "PW",
		CountryCodePY: "PY",
		CountryCodeQA: "QA",
		CountryCodeRE: "RE",
		CountryCodeRO: "RO",
		CountryCodeRS: "RS",
		CountryCodeRU: "RU",
		CountryCodeRW: "RW",
		CountryCodeSA: "SA",
		CountryCodeSB: "SB",
		CountryCodeSC: "SC",
		CountryCodeSD: "SD",
		CountryCodeSE: "SE",
		CountryCodeSG: "SG",
		CountryCodeSH: "SH",
		CountryCodeSI: "SI",
		CountryCodeSJ: "SJ",
		CountryCodeSK: "SK",
		CountryCodeSL: "SL",
		CountryCodeSM: "SM",
		CountryCodeSN: "SN",
		CountryCodeSO: "SO",
		CountryCodeSR: "SR",
		CountryCodeSS: "SS",
		CountryCodeST: "ST",
		CountryCodeSV: "SV",
		CountryCodeSX: "SX",
		CountryCodeSY: "SY",
		CountryCodeSZ: "SZ",
		CountryCodeTC: "TC",
		CountryCodeTD: "TD",
		CountryCodeTF: "TF",
		CountryCodeTG: "TG",
		CountryCodeTH: "TH",
		CountryCodeTJ: "TJ",
		CountryCodeTK: "TK",
		CountryCodeTL: "TL",
		CountryCodeTM: "TM",
		CountryCodeTN: "TN",
		CountryCodeTO: "TO",
		CountryCodeTR: "TR",
		CountryCodeTT: "TT",
		CountryCodeTV: "TV",
		CountryCodeTW: "TW",
		CountryCodeTZ: "TZ",
		CountryCodeUA: "UA",
		CountryCodeUG: "UG",
		CountryCodeUM: "UM",
		CountryCodeUS: "US",
		CountryCodeUY: "UY",
		CountryCodeUZ: "UZ",
		CountryCodeVA: "VA",
		CountryCodeVC: "VC",
		CountryCodeVE: "VE",
		CountryCodeVG: "VG",
		CountryCodeVI: "VI",
		CountryCodeVN: "VN",
		CountryCodeVU: "VU",
		CountryCodeWF: "WF",
		CountryCodeWS: "WS",
		CountryCodeYE: "YE",
		CountryCodeYT: "YT",
		CountryCodeZA: "ZA",
		CountryCodeZM: "ZM",
		CountryCodeZW: "ZW",
	}

	// SimpleStringToCountryAlpha2Code is a map from simple string to CountryAlpha2Code.
	SimpleStringToCountryAlpha2Code = map[string]CountryAlpha2Code{
		"AD": CountryCodeAD,
		"AE": CountryCodeAE,
		"AF": CountryCodeAF,
		"AG": CountryCodeAG,
		"AI": CountryCodeAI,
		"AL": CountryCodeAL,
		"AM": CountryCodeAM,
		"AO": CountryCodeAO,
		"AQ": CountryCodeAQ,
		"AR": CountryCodeAR,
		"AS": CountryCodeAS,
		"AT": CountryCodeAT,
		"AU": CountryCodeAU,
		"AW": CountryCodeAW,
		"AX": CountryCodeAX,
		"AZ": CountryCodeAZ,
		"BA": CountryCodeBA,
		"BB": CountryCodeBB,
		"BD": CountryCodeBD,
		"BE": CountryCodeBE,
		"BF": CountryCodeBF,
		"BG": CountryCodeBG,
		"BH": CountryCodeBH,
		"BI": CountryCodeBI,
		"BJ": CountryCodeBJ,
		"BL": CountryCodeBL,
		"BM": CountryCodeBM,
		"BN": CountryCodeBN,
		"BO": CountryCodeBO,
		"BQ": CountryCodeBQ,
		"BR": CountryCodeBR,
		"BS": CountryCodeBS,
		"BT": CountryCodeBT,
		"BV": CountryCodeBV,
		"BW": CountryCodeBW,
		"BY": CountryCodeBY,
		"BZ": CountryCodeBZ,
		"CA": CountryCodeCA,
		"CC": CountryCodeCC,
		"CD": CountryCodeCD,
		"CF": CountryCodeCF,
		"CG": CountryCodeCG,
		"CH": CountryCodeCH,
		"CI": CountryCodeCI,
		"CK": CountryCodeCK,
		"CL": CountryCodeCL,
		"CM": CountryCodeCM,
		"CN": CountryCodeCN,
		"CO": CountryCodeCO,
		"CR": CountryCodeCR,
		"CU": CountryCodeCU,
		"CV": CountryCodeCV,
		"CW": CountryCodeCW,
		"CX": CountryCodeCX,
		"CY": CountryCodeCY,
		"CZ": CountryCodeCZ,
		"DE": CountryCodeDE,
		"DJ": CountryCodeDJ,
		"DK": CountryCodeDK,
		"DM": CountryCodeDM,
		"DO": CountryCodeDO,
		"DZ": CountryCodeDZ,
		"EC": CountryCodeEC,
		"EE": CountryCodeEE,
		"EG": CountryCodeEG,
		"EH": CountryCodeEH,
		"ER": CountryCodeER,
		"ES": CountryCodeES,
		"ET": CountryCodeET,
		"FI": CountryCodeFI,
		"FJ": CountryCodeFJ,
		"FK": CountryCodeFK,
		"FM": CountryCodeFM,
		"FO": CountryCodeFO,
		"FR": CountryCodeFR,
		"GA": CountryCodeGA,
		"GB": CountryCodeGB,
		"GD": CountryCodeGD,
		"GE": CountryCodeGE,
		"GF": CountryCodeGF,
		"GG": CountryCodeGG,
		"GH": CountryCodeGH,
		"GI": CountryCodeGI,
		"GL": CountryCodeGL,
		"GM": CountryCodeGM,
		"GN": CountryCodeGN,
		"GP": CountryCodeGP,
		"GQ": CountryCodeGQ,
		"GR": CountryCodeGR,
		"GS": CountryCodeGS,
		"GT": CountryCodeGT,
		"GU": CountryCodeGU,
		"GW": CountryCodeGW,
		"GY": CountryCodeGY,
		"HK": CountryCodeHK,
		"HM": CountryCodeHM,
		"HN": CountryCodeHN,
		"HR": CountryCodeHR,
		"HT": CountryCodeHT,
		"HU": CountryCodeHU,
		"ID": CountryCodeID,
		"IE": CountryCodeIE,
		"IL": CountryCodeIL,
		"IM": CountryCodeIM,
		"IN": CountryCodeIN,
		"IO": CountryCodeIO,
		"IQ": CountryCodeIQ,
		"IR": CountryCodeIR,
		"IS": CountryCodeIS,
		"IT": CountryCodeIT,
		"JE": CountryCodeJE,
		"JM": CountryCodeJM,
		"JO": CountryCodeJO,
		"JP": CountryCodeJP,
		"KE": CountryCodeKE,
		"KG": CountryCodeKG,
		"KH": CountryCodeKH,
		"KI": CountryCodeKI,
		"KM": CountryCodeKM,
		"KN": CountryCodeKN,
		"KP": CountryCodeKP,
		"KR": CountryCodeKR,
		"KW": CountryCodeKW,
		"KY": CountryCodeKY,
		"KZ": CountryCodeKZ,
		"LA": CountryCodeLA,
		"LB": CountryCodeLB,
		"LC": CountryCodeLC,
		"LI": CountryCodeLI,
		"LK": CountryCodeLK,
		"LR": CountryCodeLR,
		"LS": CountryCodeLS,
		"LT": CountryCodeLT,
		"LU": CountryCodeLU,
		"LV": CountryCodeLV,
		"LY": CountryCodeLY,
		"MA": CountryCodeMA,
		"MC": CountryCodeMC,
		"MD": CountryCodeMD,
		"ME": CountryCodeME,
		"MF": CountryCodeMF,
		"MG": CountryCodeMG,
		"MH": CountryCodeMH,
		"MK": CountryCodeMK,
		"ML": CountryCodeML,
		"MM": CountryCodeMM,
		"MN": CountryCodeMN,
		"MO": CountryCodeMO,
		"MP": CountryCodeMP,
		"MQ": CountryCodeMQ,
		"MR": CountryCodeMR,
		"MS": CountryCodeMS,
		"MT": CountryCodeMT,
		"MU": CountryCodeMU,
		"MV": CountryCodeMV,
		"MW": CountryCodeMW,
		"MX": CountryCodeMX,
		"MY": CountryCodeMY,
		"MZ": CountryCodeMZ,
		"NA": CountryCodeNA,
		"NC": CountryCodeNC,
		"NE": CountryCodeNE,
		"NF": CountryCodeNF,
		"NG": CountryCodeNG,
		"NI": CountryCodeNI,
		"NL": CountryCodeNL,
		"NO": CountryCodeNO,
		"NP": CountryCodeNP,
		"NR": CountryCodeNR,
		"NU": CountryCodeNU,
		"NZ": CountryCodeNZ,
		"OM": CountryCodeOM,
		"PA": CountryCodePA,
		"PE": CountryCodePE,
		"PF": CountryCodePF,
		"PG": CountryCodePG,
		"PH": CountryCodePH,
		"PK": CountryCodePK,
		"PL": CountryCodePL,
		"PM": CountryCodePM,
		"PN": CountryCodePN,
		"PR": CountryCodePR,
		"PS": CountryCodePS,
		"PT": CountryCodePT,
		"PW": CountryCodePW,
		"PY": CountryCodePY,
		"QA": CountryCodeQA,
		"RE": CountryCodeRE,
		"RO": CountryCodeRO,
		"RS": CountryCodeRS,
		"RU": CountryCodeRU,
		"RW": CountryCodeRW,
		"SA": CountryCodeSA,
		"SB": CountryCodeSB,
		"SC": CountryCodeSC,
		"SD": CountryCodeSD,
		"SE": CountryCodeSE,
		"SG": CountryCodeSG,
		"SH": CountryCodeSH,
		"SI": CountryCodeSI,
		"SJ": CountryCodeSJ,
		"SK": CountryCodeSK,
		"SL": CountryCodeSL,
		"SM": CountryCodeSM,
		"SN": CountryCodeSN,
		"SO": CountryCodeSO,
		"SR": CountryCodeSR,
		"SS": CountryCodeSS,
		"ST": CountryCodeST,
		"SV": CountryCodeSV,
		"SX": CountryCodeSX,
		"SY": CountryCodeSY,
		"SZ": CountryCodeSZ,
		"TC": CountryCodeTC,
		"TD": CountryCodeTD,
		"TF": CountryCodeTF,
		"TG": CountryCodeTG,
		"TH": CountryCodeTH,
		"TJ": CountryCodeTJ,
		"TK": CountryCodeTK,
		"TL": CountryCodeTL,
		"TM": CountryCodeTM,
		"TN": CountryCodeTN,
		"TO": CountryCodeTO,
		"TR": CountryCodeTR,
		"TT": CountryCodeTT,
		"TV": CountryCodeTV,
		"TW": CountryCodeTW,
		"TZ": CountryCodeTZ,
		"UA": CountryCodeUA,
		"UG": CountryCodeUG,
		"UM": CountryCodeUM,
		"US": CountryCodeUS,
		"UY": CountryCodeUY,
		"UZ": CountryCodeUZ,
		"VA": CountryCodeVA,
		"VC": CountryCodeVC,
		"VE": CountryCodeVE,
		"VG": CountryCodeVG,
		"VI": CountryCodeVI,
		"VN": CountryCodeVN,
		"VU": CountryCodeVU,
		"WF": CountryCodeWF,
		"WS": CountryCodeWS,
		"YE": CountryCodeYE,
		"YT": CountryCodeYT,
		"ZA": CountryCodeZA,
		"ZM": CountryCodeZM,
		"ZW": CountryCodeZW,
	}

	// CountryAlpha3CodeToSimpleString is a map from CountryAlpha3Code to simple string.
	CountryAlpha3CodeToSimpleString = map[CountryAlpha3Code]string{
		CountryCodeABW: "ABW",
		CountryCodeAFG: "AFG",
		CountryCodeAGO: "AGO",
		CountryCodeAIA: "AIA",
		CountryCodeALA: "ALA",
		CountryCodeALB: "ALB",
		CountryCodeAND: "AND",
		CountryCodeARE: "ARE",
		CountryCodeARG: "ARG",
		CountryCodeARM: "ARM",
		CountryCodeASM: "ASM",
		CountryCodeATA: "ATA",
		CountryCodeATF: "ATF",
		CountryCodeATG: "ATG",
		CountryCodeAUS: "AUS",
		CountryCodeAUT: "AUT",
		CountryCodeAZE: "AZE",
		CountryCodeBDI: "BDI",
		CountryCodeBEL: "BEL",
		CountryCodeBEN: "BEN",
		CountryCodeBES: "BES",
		CountryCodeBFA: "BFA",
		CountryCodeBGD: "BGD",
		CountryCodeBGR: "BGR",
		CountryCodeBHR: "BHR",
		CountryCodeBHS: "BHS",
		CountryCodeBIH: "BIH",
		CountryCodeBLM: "BLM",
		CountryCodeBLR: "BLR",
		CountryCodeBLZ: "BLZ",
		CountryCodeBMU: "BMU",
		CountryCodeBOL: "BOL",
		CountryCodeBRA: "BRA",
		CountryCodeBRB: "BRB",
		CountryCodeBRN: "BRN",
		CountryCodeBTN: "BTN",
		CountryCodeBVT: "BVT",
		CountryCodeBWA: "BWA",
		CountryCodeCAF: "CAF",
		CountryCodeCAN: "CAN",
		CountryCodeCCK: "CCK",
		CountryCodeCHE: "CHE",
		CountryCodeCHL: "CHL",
		CountryCodeCHN: "CHN",
		CountryCodeCIV: "CIV",
		CountryCodeCMR: "CMR",
		CountryCodeCOD: "COD",
		CountryCodeCOG: "COG",
		CountryCodeCOK: "COK",
		CountryCodeCOL: "COL",
		CountryCodeCOM: "COM",
		CountryCodeCPV: "CPV",
		CountryCodeCRI: "CRI",
		CountryCodeCUB: "CUB",
		CountryCodeCUW: "CUW",
		CountryCodeCXR: "CXR",
		CountryCodeCYM: "CYM",
		CountryCodeCYP: "CYP",
		CountryCodeCZE: "CZE",
		CountryCodeDEU: "DEU",
		CountryCodeDJI: "DJI",
		CountryCodeDMA: "DMA",
		CountryCodeDNK: "DNK",
		CountryCodeDOM: "DOM",
		CountryCodeDZA: "DZA",
		CountryCodeECU: "ECU",
		CountryCodeEGY: "EGY",
		CountryCodeERI: "ERI",
		CountryCodeESH: "ESH",
		CountryCodeESP: "ESP",
		CountryCodeEST: "EST",
		CountryCodeETH: "ETH",
		CountryCodeFIN: "FIN",
		CountryCodeFJI: "FJI",
		CountryCodeFLK: "FLK",
		CountryCodeFRA: "FRA",
		CountryCodeFRO: "FRO",
		CountryCodeFSM: "FSM",
		CountryCodeGAB: "GAB",
		CountryCodeGBR: "GBR",
		CountryCodeGEO: "GEO",
		CountryCodeGGY: "GGY",
		CountryCodeGHA: "GHA",
		CountryCodeGIB: "GIB",
		CountryCodeGIN: "GIN",
		CountryCodeGLP: "GLP",
		CountryCodeGMB: "GMB",
		CountryCodeGNB: "GNB",
		CountryCodeGNQ: "GNQ",
		CountryCodeGRC: "GRC",
		CountryCodeGRD: "GRD",
		CountryCodeGRL: "GRL",
		CountryCodeGTM: "GTM",
		CountryCodeGUF: "GUF",
		CountryCodeGUM: "GUM",
		CountryCodeGUY: "GUY",
		CountryCodeHKG: "HKG",
		CountryCodeHMD: "HMD",
		CountryCodeHND: "HND",
		CountryCodeHRV: "HRV",
		CountryCodeHTI: "HTI",
		CountryCodeHUN: "HUN",
		CountryCodeIDN: "IDN",
		CountryCodeIMN: "IMN",
		CountryCodeIND: "IND",
		CountryCodeIOT: "IOT",
		CountryCodeIRL: "IRL",
		CountryCodeIRN: "IRN",
		CountryCodeIRQ: "IRQ",
		CountryCodeISL: "ISL",
		CountryCodeISR: "ISR",
		CountryCodeITA: "ITA",
		CountryCodeJAM: "JAM",
		CountryCodeJEY: "JEY",
		CountryCodeJOR: "JOR",
		CountryCodeJPN: "JPN",
		CountryCodeKAZ: "KAZ",
		CountryCodeKEN: "KEN",
		CountryCodeKGZ: "KGZ",
		CountryCodeKHM: "KHM",
		CountryCodeKIR: "KIR",
		CountryCodeKNA: "KNA",
		CountryCodeKOR: "KOR",
		CountryCodeKWT: "KWT",
		CountryCodeLAO: "LAO",
		CountryCodeLBN: "LBN",
		CountryCodeLBR: "LBR",
		CountryCodeLBY: "LBY",
		CountryCodeLCA: "LCA",
		CountryCodeLIE: "LIE",
		CountryCodeLKA: "LKA",
		CountryCodeLSO: "LSO",
		CountryCodeLTU: "LTU",
		CountryCodeLUX: "LUX",
		CountryCodeLVA: "LVA",
		CountryCodeMAC: "MAC",
		CountryCodeMAF: "MAF",
		CountryCodeMAR: "MAR",
		CountryCodeMCO: "MCO",
		CountryCodeMDA: "MDA",
		CountryCodeMDG: "MDG",
		CountryCodeMDV: "MDV",
		CountryCodeMEX: "MEX",
		CountryCodeMHL: "MHL",
		CountryCodeMKD: "MKD",
		CountryCodeMLI: "MLI",
		CountryCodeMLT: "MLT",
		CountryCodeMMR: "MMR",
		CountryCodeMNE: "MNE",
		CountryCodeMNG: "MNG",
		CountryCodeMNP: "MNP",
		CountryCodeMOZ: "MOZ",
		CountryCodeMRT: "MRT",
		CountryCodeMSR: "MSR",
		CountryCodeMTQ: "MTQ",
		CountryCodeMUS: "MUS",
		CountryCodeMWI: "MWI",
		CountryCodeMYS: "MYS",
		CountryCodeMYT: "MYT",
		CountryCodeNAM: "NAM",
		CountryCodeNCL: "NCL",
		CountryCodeNER: "NER",
		CountryCodeNFK: "NFK",
		CountryCodeNGA: "NGA",
		CountryCodeNIC: "NIC",
		CountryCodeNIU: "NIU",
		CountryCodeNLD: "NLD",
		CountryCodeNOR: "NOR",
		CountryCodeNPL: "NPL",
		CountryCodeNRU: "NRU",
		CountryCodeNZL: "NZL",
		CountryCodeOMN: "OMN",
		CountryCodePAK: "PAK",
		CountryCodePAN: "PAN",
		CountryCodePCN: "PCN",
		CountryCodePER: "PER",
		CountryCodePHL: "PHL",
		CountryCodePLW: "PLW",
		CountryCodePNG: "PNG",
		CountryCodePOL: "POL",
		CountryCodePRI: "PRI",
		CountryCodePRK: "PRK",
		CountryCodePRT: "PRT",
		CountryCodePRY: "PRY",
		CountryCodePSE: "PSE",
		CountryCodePYF: "PYF",
		CountryCodeQAT: "QAT",
		CountryCodeREU: "REU",
		CountryCodeROU: "ROU",
		CountryCodeRUS: "RUS",
		CountryCodeRWA: "RWA",
		CountryCodeSAU: "SAU",
		CountryCodeSDN: "SDN",
		CountryCodeSEN: "SEN",
		CountryCodeSGP: "SGP",
		CountryCodeSGS: "SGS",
		CountryCodeSHN: "SHN",
		CountryCodeSJM: "SJM",
		CountryCodeSLB: "SLB",
		CountryCodeSLE: "SLE",
		CountryCodeSLV: "SLV",
		CountryCodeSMR: "SMR",
		CountryCodeSOM: "SOM",
		CountryCodeSPM: "SPM",
		CountryCodeSRB: "SRB",
		CountryCodeSSD: "SSD",
		CountryCodeSTP: "STP",
		CountryCodeSUR: "SUR",
		CountryCodeSVK: "SVK",
		CountryCodeSVN: "SVN",
		CountryCodeSWE: "SWE",
		CountryCodeSWZ: "SWZ",
		CountryCodeSXM: "SXM",
		CountryCodeSYC: "SYC",
		CountryCodeSYR: "SYR",
		CountryCodeTCA: "TCA",
		CountryCodeTCD: "TCD",
		CountryCodeTGO: "TGO",
		CountryCodeTHA: "THA",
		CountryCodeTJK: "TJK",
		CountryCodeTKL: "TKL",
		CountryCodeTKM: "TKM",
		CountryCodeTLS: "TLS",
		CountryCodeTON: "TON",
		CountryCodeTTO: "TTO",
		CountryCodeTUN: "TUN",
		CountryCodeTUR: "TUR",
		CountryCodeTUV: "TUV",
		CountryCodeTWN: "TWN",
		CountryCodeTZA: "TZA",
		CountryCodeUGA: "UGA",
		CountryCodeUKR: "UKR",
		CountryCodeUMI: "UMI",
		CountryCodeURY: "URY",
		CountryCodeUSA: "USA",
		CountryCodeUZB: "UZB",
		CountryCodeVAT: "VAT",
		CountryCodeVCT: "VCT",
		CountryCodeVEN: "VEN",
		CountryCodeVGB: "VGB",
		CountryCodeVIR: "VIR",
		CountryCodeVNM: "VNM",
		CountryCodeVUT: "VUT",
		CountryCodeWLF: "WLF",
		CountryCodeWSM: "WSM",
		CountryCodeYEM: "YEM",
		CountryCodeZAF: "ZAF",
		CountryCodeZMB: "ZMB",
		CountryCodeZWE: "ZWE",
	}

	// SimpleStringToCountryAlpha3Code is a map from simple string to CountryAlpha3Code.
	SimpleStringToCountryAlpha3Code = map[string]CountryAlpha3Code{
		"ABW": CountryCodeABW,
		"AFG": CountryCodeAFG,
		"AGO": CountryCodeAGO,
		"AIA": CountryCodeAIA,
		"ALA": CountryCodeALA,
		"ALB": CountryCodeALB,
		"AND": CountryCodeAND,
		"ARE": CountryCodeARE,
		"ARG": CountryCodeARG,
		"ARM": CountryCodeARM,
		"ASM": CountryCodeASM,
		"ATA": CountryCodeATA,
		"ATF": CountryCodeATF,
		"ATG": CountryCodeATG,
		"AUS": CountryCodeAUS,
		"AUT": CountryCodeAUT,
		"AZE": CountryCodeAZE,
		"BDI": CountryCodeBDI,
		"BEL": CountryCodeBEL,
		"BEN": CountryCodeBEN,
		"BES": CountryCodeBES,
		"BFA": CountryCodeBFA,
		"BGD": CountryCodeBGD,
		"BGR": CountryCodeBGR,
		"BHR": CountryCodeBHR,
		"BHS": CountryCodeBHS,
		"BIH": CountryCodeBIH,
		"BLM": CountryCodeBLM,
		"BLR": CountryCodeBLR,
		"BLZ": CountryCodeBLZ,
		"BMU": CountryCodeBMU,
		"BOL": CountryCodeBOL,
		"BRA": CountryCodeBRA,
		"BRB": CountryCodeBRB,
		"BRN": CountryCodeBRN,
		"BTN": CountryCodeBTN,
		"BVT": CountryCodeBVT,
		"BWA": CountryCodeBWA,
		"CAF": CountryCodeCAF,
		"CAN": CountryCodeCAN,
		"CCK": CountryCodeCCK,
		"CHE": CountryCodeCHE,
		"CHL": CountryCodeCHL,
		"CHN": CountryCodeCHN,
		"CIV": CountryCodeCIV,
		"CMR": CountryCodeCMR,
		"COD": CountryCodeCOD,
		"COG": CountryCodeCOG,
		"COK": CountryCodeCOK,
		"COL": CountryCodeCOL,
		"COM": CountryCodeCOM,
		"CPV": CountryCodeCPV,
		"CRI": CountryCodeCRI,
		"CUB": CountryCodeCUB,
		"CUW": CountryCodeCUW,
		"CXR": CountryCodeCXR,
		"CYM": CountryCodeCYM,
		"CYP": CountryCodeCYP,
		"CZE": CountryCodeCZE,
		"DEU": CountryCodeDEU,
		"DJI": CountryCodeDJI,
		"DMA": CountryCodeDMA,
		"DNK": CountryCodeDNK,
		"DOM": CountryCodeDOM,
		"DZA": CountryCodeDZA,
		"ECU": CountryCodeECU,
		"EGY": CountryCodeEGY,
		"ERI": CountryCodeERI,
		"ESH": CountryCodeESH,
		"ESP": CountryCodeESP,
		"EST": CountryCodeEST,
		"ETH": CountryCodeETH,
		"FIN": CountryCodeFIN,
		"FJI": CountryCodeFJI,
		"FLK": CountryCodeFLK,
		"FRA": CountryCodeFRA,
		"FRO": CountryCodeFRO,
		"FSM": CountryCodeFSM,
		"GAB": CountryCodeGAB,
		"GBR": CountryCodeGBR,
		"GEO": CountryCodeGEO,
		"GGY": CountryCodeGGY,
		"GHA": CountryCodeGHA,
		"GIB": CountryCodeGIB,
		"GIN": CountryCodeGIN,
		"GLP": CountryCodeGLP,
		"GMB": CountryCodeGMB,
		"GNB": CountryCodeGNB,
		"GNQ": CountryCodeGNQ,
		"GRC": CountryCodeGRC,
		"GRD": CountryCodeGRD,
		"GRL": CountryCodeGRL,
		"GTM": CountryCodeGTM,
		"GUF": CountryCodeGUF,
		"GUM": CountryCodeGUM,
		"GUY": CountryCodeGUY,
		"HKG": CountryCodeHKG,
		"HMD": CountryCodeHMD,
		"HND": CountryCodeHND,
		"HRV": CountryCodeHRV,
		"HTI": CountryCodeHTI,
		"HUN": CountryCodeHUN,
		"IDN": CountryCodeIDN,
		"IMN": CountryCodeIMN,
		"IND": CountryCodeIND,
		"IOT": CountryCodeIOT,
		"IRL": CountryCodeIRL,
		"IRN": CountryCodeIRN,
		"IRQ": CountryCodeIRQ,
		"ISL": CountryCodeISL,
		"ISR": CountryCodeISR,
		"ITA": CountryCodeITA,
		"JAM": CountryCodeJAM,
		"JEY": CountryCodeJEY,
		"JOR": CountryCodeJOR,
		"JPN": CountryCodeJPN,
		"KAZ": CountryCodeKAZ,
		"KEN": CountryCodeKEN,
		"KGZ": CountryCodeKGZ,
		"KHM": CountryCodeKHM,
		"KIR": CountryCodeKIR,
		"KNA": CountryCodeKNA,
		"KOR": CountryCodeKOR,
		"KWT": CountryCodeKWT,
		"LAO": CountryCodeLAO,
		"LBN": CountryCodeLBN,
		"LBR": CountryCodeLBR,
		"LBY": CountryCodeLBY,
		"LCA": CountryCodeLCA,
		"LIE": CountryCodeLIE,
		"LKA": CountryCodeLKA,
		"LSO": CountryCodeLSO,
		"LTU": CountryCodeLTU,
		"LUX": CountryCodeLUX,
		"LVA": CountryCodeLVA,
		"MAC": CountryCodeMAC,
		"MAF": CountryCodeMAF,
		"MAR": CountryCodeMAR,
		"MCO": CountryCodeMCO,
		"MDA": CountryCodeMDA,
		"MDG": CountryCodeMDG,
		"MDV": CountryCodeMDV,
		"MEX": CountryCodeMEX,
		"MHL": CountryCodeMHL,
		"MKD": CountryCodeMKD,
		"MLI": CountryCodeMLI,
		"MLT": CountryCodeMLT,
		"MMR": CountryCodeMMR,
		"MNE": CountryCodeMNE,
		"MNG": CountryCodeMNG,
		"MNP": CountryCodeMNP,
		"MOZ": CountryCodeMOZ,
		"MRT": CountryCodeMRT,
		"MSR": CountryCodeMSR,
		"MTQ": CountryCodeMTQ,
		"MUS": CountryCodeMUS,
		"MWI": CountryCodeMWI,
		"MYS": CountryCodeMYS,
		"MYT": CountryCodeMYT,
		"NAM": CountryCodeNAM,
		"NCL": CountryCodeNCL,
		"NER": CountryCodeNER,
		"NFK": CountryCodeNFK,
		"NGA": CountryCodeNGA,
		"NIC": CountryCodeNIC,
		"NIU": CountryCodeNIU,
		"NLD": CountryCodeNLD,
		"NOR": CountryCodeNOR,
		"NPL": CountryCodeNPL,
		"NRU": CountryCodeNRU,
		"NZL": CountryCodeNZL,
		"OMN": CountryCodeOMN,
		"PAK": CountryCodePAK,
		"PAN": CountryCodePAN,
		"PCN": CountryCodePCN,
		"PER": CountryCodePER,
		"PHL": CountryCodePHL,
		"PLW": CountryCodePLW,
		"PNG": CountryCodePNG,
		"POL": CountryCodePOL,
		"PRI": CountryCodePRI,
		"PRK": CountryCodePRK,
		"PRT": CountryCodePRT,
		"PRY": CountryCodePRY,
		"PSE": CountryCodePSE,
		"PYF": CountryCodePYF,
		"QAT": CountryCodeQAT,
		"REU": CountryCodeREU,
		"ROU": CountryCodeROU,
		"RUS": CountryCodeRUS,
		"RWA": CountryCodeRWA,
		"SAU": CountryCodeSAU,
		"SDN": CountryCodeSDN,
		"SEN": CountryCodeSEN,
		"SGP": CountryCodeSGP,
		"SGS": CountryCodeSGS,
		"SHN": CountryCodeSHN,
		"SJM": CountryCodeSJM,
		"SLB": CountryCodeSLB,
		"SLE": CountryCodeSLE,
		"SLV": CountryCodeSLV,
		"SMR": CountryCodeSMR,
		"SOM": CountryCodeSOM,
		"SPM": CountryCodeSPM,
		"SRB": CountryCodeSRB,
		"SSD": CountryCodeSSD,
		"STP": CountryCodeSTP,
		"SUR": CountryCodeSUR,
		"SVK": CountryCodeSVK,
		"SVN": CountryCodeSVN,
		"SWE": CountryCodeSWE,
		"SWZ": CountryCodeSWZ,
		"SXM": CountryCodeSXM,
		"SYC": CountryCodeSYC,
		"SYR": CountryCodeSYR,
		"TCA": CountryCodeTCA,
		"TCD": CountryCodeTCD,
		"TGO": CountryCodeTGO,
		"THA": CountryCodeTHA,
		"TJK": CountryCodeTJK,
		"TKL": CountryCodeTKL,
		"TKM": CountryCodeTKM,
		"TLS": CountryCodeTLS,
		"TON": CountryCodeTON,
		"TTO": CountryCodeTTO,
		"TUN": CountryCodeTUN,
		"TUR": CountryCodeTUR,
		"TUV": CountryCodeTUV,
		"TWN": CountryCodeTWN,
		"TZA": CountryCodeTZA,
		"UGA": CountryCodeUGA,
		"UKR": CountryCodeUKR,
		"UMI": CountryCodeUMI,
		"URY": CountryCodeURY,
		"USA": CountryCodeUSA,
		"UZB": CountryCodeUZB,
		"VAT": CountryCodeVAT,
		"VCT": CountryCodeVCT,
		"VEN": CountryCodeVEN,
		"VGB": CountryCodeVGB,
		"VIR": CountryCodeVIR,
		"VNM": CountryCodeVNM,
		"VUT": CountryCodeVUT,
		"WLF": CountryCodeWLF,
		"WSM": CountryCodeWSM,
		"YEM": CountryCodeYEM,
		"ZAF": CountryCodeZAF,
		"ZMB": CountryCodeZMB,
		"ZWE": CountryCodeZWE,
	}

	country1 = &Country{
		Name:         "Andorra",
		Alpha_2Code:  CountryCodeAD,
		Alpha_3Code:  CountryCodeAND,
		NumericCode:  20,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country2 = &Country{
		Name:         "United Arab Emirates",
		Alpha_2Code:  CountryCodeAE,
		Alpha_3Code:  CountryCodeARE,
		NumericCode:  784,
		CurrencyCode: pbmoney.CurrencyCodeAED,
	}
	country3 = &Country{
		Name:         "Afghanistan",
		Alpha_2Code:  CountryCodeAF,
		Alpha_3Code:  CountryCodeAFG,
		NumericCode:  4,
		CurrencyCode: pbmoney.CurrencyCodeAFN,
	}
	country4 = &Country{
		Name:         "Antigua and Barbuda",
		Alpha_2Code:  CountryCodeAG,
		Alpha_3Code:  CountryCodeATG,
		NumericCode:  28,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country5 = &Country{
		Name:         "Anguilla",
		Alpha_2Code:  CountryCodeAI,
		Alpha_3Code:  CountryCodeAIA,
		NumericCode:  660,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country6 = &Country{
		Name:         "Albania",
		Alpha_2Code:  CountryCodeAL,
		Alpha_3Code:  CountryCodeALB,
		NumericCode:  8,
		CurrencyCode: pbmoney.CurrencyCodeALL,
	}
	country7 = &Country{
		Name:         "Armenia",
		Alpha_2Code:  CountryCodeAM,
		Alpha_3Code:  CountryCodeARM,
		NumericCode:  51,
		CurrencyCode: pbmoney.CurrencyCodeAMD,
	}
	country8 = &Country{
		Name:         "Angola",
		Alpha_2Code:  CountryCodeAO,
		Alpha_3Code:  CountryCodeAGO,
		NumericCode:  24,
		CurrencyCode: pbmoney.CurrencyCodeAOA,
	}
	country9 = &Country{
		Name:         "Antarctica",
		Alpha_2Code:  CountryCodeAQ,
		Alpha_3Code:  CountryCodeATA,
		NumericCode:  10,
		CurrencyCode: pbmoney.CurrencyCode_CURRENCY_CODE_NONE,
	}
	country10 = &Country{
		Name:         "Argentina",
		Alpha_2Code:  CountryCodeAR,
		Alpha_3Code:  CountryCodeARG,
		NumericCode:  32,
		CurrencyCode: pbmoney.CurrencyCodeARS,
	}
	country11 = &Country{
		Name:         "American Samoa",
		Alpha_2Code:  CountryCodeAS,
		Alpha_3Code:  CountryCodeASM,
		NumericCode:  16,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country12 = &Country{
		Name:         "Austria",
		Alpha_2Code:  CountryCodeAT,
		Alpha_3Code:  CountryCodeAUT,
		NumericCode:  40,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country13 = &Country{
		Name:         "Australia",
		Alpha_2Code:  CountryCodeAU,
		Alpha_3Code:  CountryCodeAUS,
		NumericCode:  36,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country14 = &Country{
		Name:         "Aruba",
		Alpha_2Code:  CountryCodeAW,
		Alpha_3Code:  CountryCodeABW,
		NumericCode:  533,
		CurrencyCode: pbmoney.CurrencyCodeAWG,
	}
	country15 = &Country{
		Name:         "Åland Islands",
		Alpha_2Code:  CountryCodeAX,
		Alpha_3Code:  CountryCodeALA,
		NumericCode:  248,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country16 = &Country{
		Name:         "Azerbaijan",
		Alpha_2Code:  CountryCodeAZ,
		Alpha_3Code:  CountryCodeAZE,
		NumericCode:  31,
		CurrencyCode: pbmoney.CurrencyCodeAZN,
	}
	country17 = &Country{
		Name:         "Bosnia and Herzegovina",
		Alpha_2Code:  CountryCodeBA,
		Alpha_3Code:  CountryCodeBIH,
		NumericCode:  70,
		CurrencyCode: pbmoney.CurrencyCodeBAM,
	}
	country18 = &Country{
		Name:         "Barbados",
		Alpha_2Code:  CountryCodeBB,
		Alpha_3Code:  CountryCodeBRB,
		NumericCode:  52,
		CurrencyCode: pbmoney.CurrencyCodeBBD,
	}
	country19 = &Country{
		Name:         "Bangladesh",
		Alpha_2Code:  CountryCodeBD,
		Alpha_3Code:  CountryCodeBGD,
		NumericCode:  50,
		CurrencyCode: pbmoney.CurrencyCodeBDT,
	}
	country20 = &Country{
		Name:         "Belgium",
		Alpha_2Code:  CountryCodeBE,
		Alpha_3Code:  CountryCodeBEL,
		NumericCode:  56,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country21 = &Country{
		Name:         "Burkina Faso",
		Alpha_2Code:  CountryCodeBF,
		Alpha_3Code:  CountryCodeBFA,
		NumericCode:  854,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country22 = &Country{
		Name:         "Bulgaria",
		Alpha_2Code:  CountryCodeBG,
		Alpha_3Code:  CountryCodeBGR,
		NumericCode:  100,
		CurrencyCode: pbmoney.CurrencyCodeBGN,
	}
	country23 = &Country{
		Name:         "Bahrain",
		Alpha_2Code:  CountryCodeBH,
		Alpha_3Code:  CountryCodeBHR,
		NumericCode:  48,
		CurrencyCode: pbmoney.CurrencyCodeBHD,
	}
	country24 = &Country{
		Name:         "Burundi",
		Alpha_2Code:  CountryCodeBI,
		Alpha_3Code:  CountryCodeBDI,
		NumericCode:  108,
		CurrencyCode: pbmoney.CurrencyCodeBIF,
	}
	country25 = &Country{
		Name:         "Benin",
		Alpha_2Code:  CountryCodeBJ,
		Alpha_3Code:  CountryCodeBEN,
		NumericCode:  204,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country26 = &Country{
		Name:         "Saint Barthélemy",
		Alpha_2Code:  CountryCodeBL,
		Alpha_3Code:  CountryCodeBLM,
		NumericCode:  652,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country27 = &Country{
		Name:         "Bermuda",
		Alpha_2Code:  CountryCodeBM,
		Alpha_3Code:  CountryCodeBMU,
		NumericCode:  60,
		CurrencyCode: pbmoney.CurrencyCodeBMD,
	}
	country28 = &Country{
		Name:         "Brunei Darussalam",
		Alpha_2Code:  CountryCodeBN,
		Alpha_3Code:  CountryCodeBRN,
		NumericCode:  96,
		CurrencyCode: pbmoney.CurrencyCodeBND,
	}
	country29 = &Country{
		Name:         "Bolivia, Plurinational State of",
		Alpha_2Code:  CountryCodeBO,
		Alpha_3Code:  CountryCodeBOL,
		NumericCode:  68,
		CurrencyCode: pbmoney.CurrencyCodeBOB,
	}
	country30 = &Country{
		Name:         "Bonaire, Sint Eustatius and Saba",
		Alpha_2Code:  CountryCodeBQ,
		Alpha_3Code:  CountryCodeBES,
		NumericCode:  535,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country31 = &Country{
		Name:         "Brazil",
		Alpha_2Code:  CountryCodeBR,
		Alpha_3Code:  CountryCodeBRA,
		NumericCode:  76,
		CurrencyCode: pbmoney.CurrencyCodeBRL,
	}
	country32 = &Country{
		Name:         "Bahamas",
		Alpha_2Code:  CountryCodeBS,
		Alpha_3Code:  CountryCodeBHS,
		NumericCode:  44,
		CurrencyCode: pbmoney.CurrencyCodeBSD,
	}
	country33 = &Country{
		Name:         "Bhutan",
		Alpha_2Code:  CountryCodeBT,
		Alpha_3Code:  CountryCodeBTN,
		NumericCode:  64,
		CurrencyCode: pbmoney.CurrencyCodeINR,
	}
	country34 = &Country{
		Name:         "Bouvet Island",
		Alpha_2Code:  CountryCodeBV,
		Alpha_3Code:  CountryCodeBVT,
		NumericCode:  74,
		CurrencyCode: pbmoney.CurrencyCodeNOK,
	}
	country35 = &Country{
		Name:         "Botswana",
		Alpha_2Code:  CountryCodeBW,
		Alpha_3Code:  CountryCodeBWA,
		NumericCode:  72,
		CurrencyCode: pbmoney.CurrencyCodeBWP,
	}
	country36 = &Country{
		Name:         "Belarus",
		Alpha_2Code:  CountryCodeBY,
		Alpha_3Code:  CountryCodeBLR,
		NumericCode:  112,
		CurrencyCode: pbmoney.CurrencyCodeBYR,
	}
	country37 = &Country{
		Name:         "Belize",
		Alpha_2Code:  CountryCodeBZ,
		Alpha_3Code:  CountryCodeBLZ,
		NumericCode:  84,
		CurrencyCode: pbmoney.CurrencyCodeBZD,
	}
	country38 = &Country{
		Name:         "Canada",
		Alpha_2Code:  CountryCodeCA,
		Alpha_3Code:  CountryCodeCAN,
		NumericCode:  124,
		CurrencyCode: pbmoney.CurrencyCodeCAD,
	}
	country39 = &Country{
		Name:         "Cocos (Keeling) Islands",
		Alpha_2Code:  CountryCodeCC,
		Alpha_3Code:  CountryCodeCCK,
		NumericCode:  166,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country40 = &Country{
		Name:         "Congo, the Democratic Republic of the",
		Alpha_2Code:  CountryCodeCD,
		Alpha_3Code:  CountryCodeCOD,
		NumericCode:  180,
		CurrencyCode: pbmoney.CurrencyCode_CURRENCY_CODE_NONE,
	}
	country41 = &Country{
		Name:         "Central African Republic",
		Alpha_2Code:  CountryCodeCF,
		Alpha_3Code:  CountryCodeCAF,
		NumericCode:  140,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country42 = &Country{
		Name:         "Congo",
		Alpha_2Code:  CountryCodeCG,
		Alpha_3Code:  CountryCodeCOG,
		NumericCode:  178,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country43 = &Country{
		Name:         "Switzerland",
		Alpha_2Code:  CountryCodeCH,
		Alpha_3Code:  CountryCodeCHE,
		NumericCode:  756,
		CurrencyCode: pbmoney.CurrencyCodeCHF,
	}
	country44 = &Country{
		Name:         "Côte d'Ivoire",
		Alpha_2Code:  CountryCodeCI,
		Alpha_3Code:  CountryCodeCIV,
		NumericCode:  384,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country45 = &Country{
		Name:         "Cook Islands",
		Alpha_2Code:  CountryCodeCK,
		Alpha_3Code:  CountryCodeCOK,
		NumericCode:  184,
		CurrencyCode: pbmoney.CurrencyCodeNZD,
	}
	country46 = &Country{
		Name:         "Chile",
		Alpha_2Code:  CountryCodeCL,
		Alpha_3Code:  CountryCodeCHL,
		NumericCode:  152,
		CurrencyCode: pbmoney.CurrencyCodeCLP,
	}
	country47 = &Country{
		Name:         "Cameroon",
		Alpha_2Code:  CountryCodeCM,
		Alpha_3Code:  CountryCodeCMR,
		NumericCode:  120,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country48 = &Country{
		Name:         "China",
		Alpha_2Code:  CountryCodeCN,
		Alpha_3Code:  CountryCodeCHN,
		NumericCode:  156,
		CurrencyCode: pbmoney.CurrencyCodeCNY,
	}
	country49 = &Country{
		Name:         "Colombia",
		Alpha_2Code:  CountryCodeCO,
		Alpha_3Code:  CountryCodeCOL,
		NumericCode:  170,
		CurrencyCode: pbmoney.CurrencyCodeCOP,
	}
	country50 = &Country{
		Name:         "Costa Rica",
		Alpha_2Code:  CountryCodeCR,
		Alpha_3Code:  CountryCodeCRI,
		NumericCode:  188,
		CurrencyCode: pbmoney.CurrencyCodeCRC,
	}
	country51 = &Country{
		Name:         "Cuba",
		Alpha_2Code:  CountryCodeCU,
		Alpha_3Code:  CountryCodeCUB,
		NumericCode:  192,
		CurrencyCode: pbmoney.CurrencyCodeCUP,
	}
	country52 = &Country{
		Name:         "Cape Verde",
		Alpha_2Code:  CountryCodeCV,
		Alpha_3Code:  CountryCodeCPV,
		NumericCode:  132,
		CurrencyCode: pbmoney.CurrencyCodeCVE,
	}
	country53 = &Country{
		Name:         "Curaçao",
		Alpha_2Code:  CountryCodeCW,
		Alpha_3Code:  CountryCodeCUW,
		NumericCode:  531,
		CurrencyCode: pbmoney.CurrencyCodeANG,
	}
	country54 = &Country{
		Name:         "Christmas Island",
		Alpha_2Code:  CountryCodeCX,
		Alpha_3Code:  CountryCodeCXR,
		NumericCode:  162,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country55 = &Country{
		Name:         "Cyprus",
		Alpha_2Code:  CountryCodeCY,
		Alpha_3Code:  CountryCodeCYP,
		NumericCode:  196,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country56 = &Country{
		Name:         "Czech Republic",
		Alpha_2Code:  CountryCodeCZ,
		Alpha_3Code:  CountryCodeCZE,
		NumericCode:  203,
		CurrencyCode: pbmoney.CurrencyCodeCZK,
	}
	country57 = &Country{
		Name:         "Germany",
		Alpha_2Code:  CountryCodeDE,
		Alpha_3Code:  CountryCodeDEU,
		NumericCode:  276,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country58 = &Country{
		Name:         "Djibouti",
		Alpha_2Code:  CountryCodeDJ,
		Alpha_3Code:  CountryCodeDJI,
		NumericCode:  262,
		CurrencyCode: pbmoney.CurrencyCodeDJF,
	}
	country59 = &Country{
		Name:         "Denmark",
		Alpha_2Code:  CountryCodeDK,
		Alpha_3Code:  CountryCodeDNK,
		NumericCode:  208,
		CurrencyCode: pbmoney.CurrencyCodeDKK,
	}
	country60 = &Country{
		Name:         "Dominica",
		Alpha_2Code:  CountryCodeDM,
		Alpha_3Code:  CountryCodeDMA,
		NumericCode:  212,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country61 = &Country{
		Name:         "Dominican Republic",
		Alpha_2Code:  CountryCodeDO,
		Alpha_3Code:  CountryCodeDOM,
		NumericCode:  214,
		CurrencyCode: pbmoney.CurrencyCodeDOP,
	}
	country62 = &Country{
		Name:         "Algeria",
		Alpha_2Code:  CountryCodeDZ,
		Alpha_3Code:  CountryCodeDZA,
		NumericCode:  12,
		CurrencyCode: pbmoney.CurrencyCodeDZD,
	}
	country63 = &Country{
		Name:         "Ecuador",
		Alpha_2Code:  CountryCodeEC,
		Alpha_3Code:  CountryCodeECU,
		NumericCode:  218,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country64 = &Country{
		Name:         "Estonia",
		Alpha_2Code:  CountryCodeEE,
		Alpha_3Code:  CountryCodeEST,
		NumericCode:  233,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country65 = &Country{
		Name:         "Egypt",
		Alpha_2Code:  CountryCodeEG,
		Alpha_3Code:  CountryCodeEGY,
		NumericCode:  818,
		CurrencyCode: pbmoney.CurrencyCodeEGP,
	}
	country66 = &Country{
		Name:         "Western Sahara",
		Alpha_2Code:  CountryCodeEH,
		Alpha_3Code:  CountryCodeESH,
		NumericCode:  732,
		CurrencyCode: pbmoney.CurrencyCodeMAD,
	}
	country67 = &Country{
		Name:         "Eritrea",
		Alpha_2Code:  CountryCodeER,
		Alpha_3Code:  CountryCodeERI,
		NumericCode:  232,
		CurrencyCode: pbmoney.CurrencyCodeERN,
	}
	country68 = &Country{
		Name:         "Spain",
		Alpha_2Code:  CountryCodeES,
		Alpha_3Code:  CountryCodeESP,
		NumericCode:  724,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country69 = &Country{
		Name:         "Ethiopia",
		Alpha_2Code:  CountryCodeET,
		Alpha_3Code:  CountryCodeETH,
		NumericCode:  231,
		CurrencyCode: pbmoney.CurrencyCodeETB,
	}
	country70 = &Country{
		Name:         "Finland",
		Alpha_2Code:  CountryCodeFI,
		Alpha_3Code:  CountryCodeFIN,
		NumericCode:  246,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country71 = &Country{
		Name:         "Fiji",
		Alpha_2Code:  CountryCodeFJ,
		Alpha_3Code:  CountryCodeFJI,
		NumericCode:  242,
		CurrencyCode: pbmoney.CurrencyCodeFJD,
	}
	country72 = &Country{
		Name:         "Falkland Islands (Malvinas)",
		Alpha_2Code:  CountryCodeFK,
		Alpha_3Code:  CountryCodeFLK,
		NumericCode:  238,
		CurrencyCode: pbmoney.CurrencyCodeFKP,
	}
	country73 = &Country{
		Name:         "Micronesia, Federated States of",
		Alpha_2Code:  CountryCodeFM,
		Alpha_3Code:  CountryCodeFSM,
		NumericCode:  583,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country74 = &Country{
		Name:         "Faroe Islands",
		Alpha_2Code:  CountryCodeFO,
		Alpha_3Code:  CountryCodeFRO,
		NumericCode:  234,
		CurrencyCode: pbmoney.CurrencyCodeDKK,
	}
	country75 = &Country{
		Name:         "France",
		Alpha_2Code:  CountryCodeFR,
		Alpha_3Code:  CountryCodeFRA,
		NumericCode:  250,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country76 = &Country{
		Name:         "Gabon",
		Alpha_2Code:  CountryCodeGA,
		Alpha_3Code:  CountryCodeGAB,
		NumericCode:  266,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country77 = &Country{
		Name:         "United Kingdom",
		Alpha_2Code:  CountryCodeGB,
		Alpha_3Code:  CountryCodeGBR,
		NumericCode:  826,
		CurrencyCode: pbmoney.CurrencyCodeGBP,
	}
	country78 = &Country{
		Name:         "Grenada",
		Alpha_2Code:  CountryCodeGD,
		Alpha_3Code:  CountryCodeGRD,
		NumericCode:  308,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country79 = &Country{
		Name:         "Georgia",
		Alpha_2Code:  CountryCodeGE,
		Alpha_3Code:  CountryCodeGEO,
		NumericCode:  268,
		CurrencyCode: pbmoney.CurrencyCodeGEL,
	}
	country80 = &Country{
		Name:         "French Guiana",
		Alpha_2Code:  CountryCodeGF,
		Alpha_3Code:  CountryCodeGUF,
		NumericCode:  254,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country81 = &Country{
		Name:         "Guernsey",
		Alpha_2Code:  CountryCodeGG,
		Alpha_3Code:  CountryCodeGGY,
		NumericCode:  831,
		CurrencyCode: pbmoney.CurrencyCodeGBP,
	}
	country82 = &Country{
		Name:         "Ghana",
		Alpha_2Code:  CountryCodeGH,
		Alpha_3Code:  CountryCodeGHA,
		NumericCode:  288,
		CurrencyCode: pbmoney.CurrencyCodeGHS,
	}
	country83 = &Country{
		Name:         "Gibraltar",
		Alpha_2Code:  CountryCodeGI,
		Alpha_3Code:  CountryCodeGIB,
		NumericCode:  292,
		CurrencyCode: pbmoney.CurrencyCodeGIP,
	}
	country84 = &Country{
		Name:         "Greenland",
		Alpha_2Code:  CountryCodeGL,
		Alpha_3Code:  CountryCodeGRL,
		NumericCode:  304,
		CurrencyCode: pbmoney.CurrencyCodeDKK,
	}
	country85 = &Country{
		Name:         "Gambia",
		Alpha_2Code:  CountryCodeGM,
		Alpha_3Code:  CountryCodeGMB,
		NumericCode:  270,
		CurrencyCode: pbmoney.CurrencyCodeGMD,
	}
	country86 = &Country{
		Name:         "Guinea",
		Alpha_2Code:  CountryCodeGN,
		Alpha_3Code:  CountryCodeGIN,
		NumericCode:  324,
		CurrencyCode: pbmoney.CurrencyCodeGNF,
	}
	country87 = &Country{
		Name:         "Guadeloupe",
		Alpha_2Code:  CountryCodeGP,
		Alpha_3Code:  CountryCodeGLP,
		NumericCode:  312,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country88 = &Country{
		Name:         "Equatorial Guinea",
		Alpha_2Code:  CountryCodeGQ,
		Alpha_3Code:  CountryCodeGNQ,
		NumericCode:  226,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country89 = &Country{
		Name:         "Greece",
		Alpha_2Code:  CountryCodeGR,
		Alpha_3Code:  CountryCodeGRC,
		NumericCode:  300,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country90 = &Country{
		Name:         "South Georgia and the South Sandwich Islands",
		Alpha_2Code:  CountryCodeGS,
		Alpha_3Code:  CountryCodeSGS,
		NumericCode:  239,
		CurrencyCode: pbmoney.CurrencyCode_CURRENCY_CODE_NONE,
	}
	country91 = &Country{
		Name:         "Guatemala",
		Alpha_2Code:  CountryCodeGT,
		Alpha_3Code:  CountryCodeGTM,
		NumericCode:  320,
		CurrencyCode: pbmoney.CurrencyCodeGTQ,
	}
	country92 = &Country{
		Name:         "Guam",
		Alpha_2Code:  CountryCodeGU,
		Alpha_3Code:  CountryCodeGUM,
		NumericCode:  316,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country93 = &Country{
		Name:         "Guinea-Bissau",
		Alpha_2Code:  CountryCodeGW,
		Alpha_3Code:  CountryCodeGNB,
		NumericCode:  624,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country94 = &Country{
		Name:         "Guyana",
		Alpha_2Code:  CountryCodeGY,
		Alpha_3Code:  CountryCodeGUY,
		NumericCode:  328,
		CurrencyCode: pbmoney.CurrencyCodeGYD,
	}
	country95 = &Country{
		Name:         "Hong Kong",
		Alpha_2Code:  CountryCodeHK,
		Alpha_3Code:  CountryCodeHKG,
		NumericCode:  344,
		CurrencyCode: pbmoney.CurrencyCodeHKD,
	}
	country96 = &Country{
		Name:         "Heard Island and McDonald Islands",
		Alpha_2Code:  CountryCodeHM,
		Alpha_3Code:  CountryCodeHMD,
		NumericCode:  334,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country97 = &Country{
		Name:         "Honduras",
		Alpha_2Code:  CountryCodeHN,
		Alpha_3Code:  CountryCodeHND,
		NumericCode:  340,
		CurrencyCode: pbmoney.CurrencyCodeHNL,
	}
	country98 = &Country{
		Name:         "Croatia",
		Alpha_2Code:  CountryCodeHR,
		Alpha_3Code:  CountryCodeHRV,
		NumericCode:  191,
		CurrencyCode: pbmoney.CurrencyCodeHRK,
	}
	country99 = &Country{
		Name:         "Haiti",
		Alpha_2Code:  CountryCodeHT,
		Alpha_3Code:  CountryCodeHTI,
		NumericCode:  332,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country100 = &Country{
		Name:         "Hungary",
		Alpha_2Code:  CountryCodeHU,
		Alpha_3Code:  CountryCodeHUN,
		NumericCode:  348,
		CurrencyCode: pbmoney.CurrencyCodeHUF,
	}
	country101 = &Country{
		Name:         "Indonesia",
		Alpha_2Code:  CountryCodeID,
		Alpha_3Code:  CountryCodeIDN,
		NumericCode:  360,
		CurrencyCode: pbmoney.CurrencyCodeIDR,
	}
	country102 = &Country{
		Name:         "Ireland",
		Alpha_2Code:  CountryCodeIE,
		Alpha_3Code:  CountryCodeIRL,
		NumericCode:  372,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country103 = &Country{
		Name:         "Israel",
		Alpha_2Code:  CountryCodeIL,
		Alpha_3Code:  CountryCodeISR,
		NumericCode:  376,
		CurrencyCode: pbmoney.CurrencyCodeILS,
	}
	country104 = &Country{
		Name:         "Isle of Man",
		Alpha_2Code:  CountryCodeIM,
		Alpha_3Code:  CountryCodeIMN,
		NumericCode:  833,
		CurrencyCode: pbmoney.CurrencyCodeGBP,
	}
	country105 = &Country{
		Name:         "India",
		Alpha_2Code:  CountryCodeIN,
		Alpha_3Code:  CountryCodeIND,
		NumericCode:  356,
		CurrencyCode: pbmoney.CurrencyCodeINR,
	}
	country106 = &Country{
		Name:         "British Indian Ocean Territory",
		Alpha_2Code:  CountryCodeIO,
		Alpha_3Code:  CountryCodeIOT,
		NumericCode:  86,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country107 = &Country{
		Name:         "Iraq",
		Alpha_2Code:  CountryCodeIQ,
		Alpha_3Code:  CountryCodeIRQ,
		NumericCode:  368,
		CurrencyCode: pbmoney.CurrencyCodeIQD,
	}
	country108 = &Country{
		Name:         "Iran, Islamic Republic of",
		Alpha_2Code:  CountryCodeIR,
		Alpha_3Code:  CountryCodeIRN,
		NumericCode:  364,
		CurrencyCode: pbmoney.CurrencyCodeIRR,
	}
	country109 = &Country{
		Name:         "Iceland",
		Alpha_2Code:  CountryCodeIS,
		Alpha_3Code:  CountryCodeISL,
		NumericCode:  352,
		CurrencyCode: pbmoney.CurrencyCodeISK,
	}
	country110 = &Country{
		Name:         "Italy",
		Alpha_2Code:  CountryCodeIT,
		Alpha_3Code:  CountryCodeITA,
		NumericCode:  380,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country111 = &Country{
		Name:         "Jersey",
		Alpha_2Code:  CountryCodeJE,
		Alpha_3Code:  CountryCodeJEY,
		NumericCode:  832,
		CurrencyCode: pbmoney.CurrencyCodeGBP,
	}
	country112 = &Country{
		Name:         "Jamaica",
		Alpha_2Code:  CountryCodeJM,
		Alpha_3Code:  CountryCodeJAM,
		NumericCode:  388,
		CurrencyCode: pbmoney.CurrencyCodeJMD,
	}
	country113 = &Country{
		Name:         "Jordan",
		Alpha_2Code:  CountryCodeJO,
		Alpha_3Code:  CountryCodeJOR,
		NumericCode:  400,
		CurrencyCode: pbmoney.CurrencyCodeJOD,
	}
	country114 = &Country{
		Name:         "Japan",
		Alpha_2Code:  CountryCodeJP,
		Alpha_3Code:  CountryCodeJPN,
		NumericCode:  392,
		CurrencyCode: pbmoney.CurrencyCodeJPY,
	}
	country115 = &Country{
		Name:         "Kenya",
		Alpha_2Code:  CountryCodeKE,
		Alpha_3Code:  CountryCodeKEN,
		NumericCode:  404,
		CurrencyCode: pbmoney.CurrencyCodeKES,
	}
	country116 = &Country{
		Name:         "Kyrgyzstan",
		Alpha_2Code:  CountryCodeKG,
		Alpha_3Code:  CountryCodeKGZ,
		NumericCode:  417,
		CurrencyCode: pbmoney.CurrencyCodeKGS,
	}
	country117 = &Country{
		Name:         "Cambodia",
		Alpha_2Code:  CountryCodeKH,
		Alpha_3Code:  CountryCodeKHM,
		NumericCode:  116,
		CurrencyCode: pbmoney.CurrencyCodeKHR,
	}
	country118 = &Country{
		Name:         "Kiribati",
		Alpha_2Code:  CountryCodeKI,
		Alpha_3Code:  CountryCodeKIR,
		NumericCode:  296,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country119 = &Country{
		Name:         "Comoros",
		Alpha_2Code:  CountryCodeKM,
		Alpha_3Code:  CountryCodeCOM,
		NumericCode:  174,
		CurrencyCode: pbmoney.CurrencyCodeKMF,
	}
	country120 = &Country{
		Name:         "Saint Kitts and Nevis",
		Alpha_2Code:  CountryCodeKN,
		Alpha_3Code:  CountryCodeKNA,
		NumericCode:  659,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country121 = &Country{
		Name:         "Korea, Democratic People's Republic of",
		Alpha_2Code:  CountryCodeKP,
		Alpha_3Code:  CountryCodePRK,
		NumericCode:  408,
		CurrencyCode: pbmoney.CurrencyCodeKPW,
	}
	country122 = &Country{
		Name:         "Korea, Republic of",
		Alpha_2Code:  CountryCodeKR,
		Alpha_3Code:  CountryCodeKOR,
		NumericCode:  410,
		CurrencyCode: pbmoney.CurrencyCodeKRW,
	}
	country123 = &Country{
		Name:         "Kuwait",
		Alpha_2Code:  CountryCodeKW,
		Alpha_3Code:  CountryCodeKWT,
		NumericCode:  414,
		CurrencyCode: pbmoney.CurrencyCodeKWD,
	}
	country124 = &Country{
		Name:         "Cayman Islands",
		Alpha_2Code:  CountryCodeKY,
		Alpha_3Code:  CountryCodeCYM,
		NumericCode:  136,
		CurrencyCode: pbmoney.CurrencyCodeKYD,
	}
	country125 = &Country{
		Name:         "Kazakhstan",
		Alpha_2Code:  CountryCodeKZ,
		Alpha_3Code:  CountryCodeKAZ,
		NumericCode:  398,
		CurrencyCode: pbmoney.CurrencyCodeKZT,
	}
	country126 = &Country{
		Name:         "Lao People's Democratic Republic",
		Alpha_2Code:  CountryCodeLA,
		Alpha_3Code:  CountryCodeLAO,
		NumericCode:  418,
		CurrencyCode: pbmoney.CurrencyCodeLAK,
	}
	country127 = &Country{
		Name:         "Lebanon",
		Alpha_2Code:  CountryCodeLB,
		Alpha_3Code:  CountryCodeLBN,
		NumericCode:  422,
		CurrencyCode: pbmoney.CurrencyCodeLBP,
	}
	country128 = &Country{
		Name:         "Saint Lucia",
		Alpha_2Code:  CountryCodeLC,
		Alpha_3Code:  CountryCodeLCA,
		NumericCode:  662,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country129 = &Country{
		Name:         "Liechtenstein",
		Alpha_2Code:  CountryCodeLI,
		Alpha_3Code:  CountryCodeLIE,
		NumericCode:  438,
		CurrencyCode: pbmoney.CurrencyCodeCHF,
	}
	country130 = &Country{
		Name:         "Sri Lanka",
		Alpha_2Code:  CountryCodeLK,
		Alpha_3Code:  CountryCodeLKA,
		NumericCode:  144,
		CurrencyCode: pbmoney.CurrencyCodeLKR,
	}
	country131 = &Country{
		Name:         "Liberia",
		Alpha_2Code:  CountryCodeLR,
		Alpha_3Code:  CountryCodeLBR,
		NumericCode:  430,
		CurrencyCode: pbmoney.CurrencyCodeLRD,
	}
	country132 = &Country{
		Name:         "Lesotho",
		Alpha_2Code:  CountryCodeLS,
		Alpha_3Code:  CountryCodeLSO,
		NumericCode:  426,
		CurrencyCode: pbmoney.CurrencyCodeZAR,
	}
	country133 = &Country{
		Name:         "Lithuania",
		Alpha_2Code:  CountryCodeLT,
		Alpha_3Code:  CountryCodeLTU,
		NumericCode:  440,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country134 = &Country{
		Name:         "Luxembourg",
		Alpha_2Code:  CountryCodeLU,
		Alpha_3Code:  CountryCodeLUX,
		NumericCode:  442,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country135 = &Country{
		Name:         "Latvia",
		Alpha_2Code:  CountryCodeLV,
		Alpha_3Code:  CountryCodeLVA,
		NumericCode:  428,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country136 = &Country{
		Name:         "Libya",
		Alpha_2Code:  CountryCodeLY,
		Alpha_3Code:  CountryCodeLBY,
		NumericCode:  434,
		CurrencyCode: pbmoney.CurrencyCodeLYD,
	}
	country137 = &Country{
		Name:         "Morocco",
		Alpha_2Code:  CountryCodeMA,
		Alpha_3Code:  CountryCodeMAR,
		NumericCode:  504,
		CurrencyCode: pbmoney.CurrencyCodeMAD,
	}
	country138 = &Country{
		Name:         "Monaco",
		Alpha_2Code:  CountryCodeMC,
		Alpha_3Code:  CountryCodeMCO,
		NumericCode:  492,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country139 = &Country{
		Name:         "Moldova, Republic of",
		Alpha_2Code:  CountryCodeMD,
		Alpha_3Code:  CountryCodeMDA,
		NumericCode:  498,
		CurrencyCode: pbmoney.CurrencyCodeMDL,
	}
	country140 = &Country{
		Name:         "Montenegro",
		Alpha_2Code:  CountryCodeME,
		Alpha_3Code:  CountryCodeMNE,
		NumericCode:  499,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country141 = &Country{
		Name:         "Saint Martin (French part)",
		Alpha_2Code:  CountryCodeMF,
		Alpha_3Code:  CountryCodeMAF,
		NumericCode:  663,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country142 = &Country{
		Name:         "Madagascar",
		Alpha_2Code:  CountryCodeMG,
		Alpha_3Code:  CountryCodeMDG,
		NumericCode:  450,
		CurrencyCode: pbmoney.CurrencyCodeMGA,
	}
	country143 = &Country{
		Name:         "Marshall Islands",
		Alpha_2Code:  CountryCodeMH,
		Alpha_3Code:  CountryCodeMHL,
		NumericCode:  584,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country144 = &Country{
		Name:         "Macedonia, the Former Yugoslav Republic of",
		Alpha_2Code:  CountryCodeMK,
		Alpha_3Code:  CountryCodeMKD,
		NumericCode:  807,
		CurrencyCode: pbmoney.CurrencyCodeMKD,
	}
	country145 = &Country{
		Name:         "Mali",
		Alpha_2Code:  CountryCodeML,
		Alpha_3Code:  CountryCodeMLI,
		NumericCode:  466,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country146 = &Country{
		Name:         "Myanmar",
		Alpha_2Code:  CountryCodeMM,
		Alpha_3Code:  CountryCodeMMR,
		NumericCode:  104,
		CurrencyCode: pbmoney.CurrencyCodeMMK,
	}
	country147 = &Country{
		Name:         "Mongolia",
		Alpha_2Code:  CountryCodeMN,
		Alpha_3Code:  CountryCodeMNG,
		NumericCode:  496,
		CurrencyCode: pbmoney.CurrencyCodeMNT,
	}
	country148 = &Country{
		Name:         "Macao",
		Alpha_2Code:  CountryCodeMO,
		Alpha_3Code:  CountryCodeMAC,
		NumericCode:  446,
		CurrencyCode: pbmoney.CurrencyCodeMOP,
	}
	country149 = &Country{
		Name:         "Northern Mariana Islands",
		Alpha_2Code:  CountryCodeMP,
		Alpha_3Code:  CountryCodeMNP,
		NumericCode:  580,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country150 = &Country{
		Name:         "Martinique",
		Alpha_2Code:  CountryCodeMQ,
		Alpha_3Code:  CountryCodeMTQ,
		NumericCode:  474,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country151 = &Country{
		Name:         "Mauritania",
		Alpha_2Code:  CountryCodeMR,
		Alpha_3Code:  CountryCodeMRT,
		NumericCode:  478,
		CurrencyCode: pbmoney.CurrencyCodeMRO,
	}
	country152 = &Country{
		Name:         "Montserrat",
		Alpha_2Code:  CountryCodeMS,
		Alpha_3Code:  CountryCodeMSR,
		NumericCode:  500,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country153 = &Country{
		Name:         "Malta",
		Alpha_2Code:  CountryCodeMT,
		Alpha_3Code:  CountryCodeMLT,
		NumericCode:  470,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country154 = &Country{
		Name:         "Mauritius",
		Alpha_2Code:  CountryCodeMU,
		Alpha_3Code:  CountryCodeMUS,
		NumericCode:  480,
		CurrencyCode: pbmoney.CurrencyCodeMUR,
	}
	country155 = &Country{
		Name:         "Maldives",
		Alpha_2Code:  CountryCodeMV,
		Alpha_3Code:  CountryCodeMDV,
		NumericCode:  462,
		CurrencyCode: pbmoney.CurrencyCodeMVR,
	}
	country156 = &Country{
		Name:         "Malawi",
		Alpha_2Code:  CountryCodeMW,
		Alpha_3Code:  CountryCodeMWI,
		NumericCode:  454,
		CurrencyCode: pbmoney.CurrencyCodeMWK,
	}
	country157 = &Country{
		Name:         "Mexico",
		Alpha_2Code:  CountryCodeMX,
		Alpha_3Code:  CountryCodeMEX,
		NumericCode:  484,
		CurrencyCode: pbmoney.CurrencyCodeMXN,
	}
	country158 = &Country{
		Name:         "Malaysia",
		Alpha_2Code:  CountryCodeMY,
		Alpha_3Code:  CountryCodeMYS,
		NumericCode:  458,
		CurrencyCode: pbmoney.CurrencyCodeMYR,
	}
	country159 = &Country{
		Name:         "Mozambique",
		Alpha_2Code:  CountryCodeMZ,
		Alpha_3Code:  CountryCodeMOZ,
		NumericCode:  508,
		CurrencyCode: pbmoney.CurrencyCodeMZN,
	}
	country160 = &Country{
		Name:         "Namibia",
		Alpha_2Code:  CountryCodeNA,
		Alpha_3Code:  CountryCodeNAM,
		NumericCode:  516,
		CurrencyCode: pbmoney.CurrencyCodeZAR,
	}
	country161 = &Country{
		Name:         "New Caledonia",
		Alpha_2Code:  CountryCodeNC,
		Alpha_3Code:  CountryCodeNCL,
		NumericCode:  540,
		CurrencyCode: pbmoney.CurrencyCodeXPF,
	}
	country162 = &Country{
		Name:         "Niger",
		Alpha_2Code:  CountryCodeNE,
		Alpha_3Code:  CountryCodeNER,
		NumericCode:  562,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country163 = &Country{
		Name:         "Norfolk Island",
		Alpha_2Code:  CountryCodeNF,
		Alpha_3Code:  CountryCodeNFK,
		NumericCode:  574,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country164 = &Country{
		Name:         "Nigeria",
		Alpha_2Code:  CountryCodeNG,
		Alpha_3Code:  CountryCodeNGA,
		NumericCode:  566,
		CurrencyCode: pbmoney.CurrencyCodeNGN,
	}
	country165 = &Country{
		Name:         "Nicaragua",
		Alpha_2Code:  CountryCodeNI,
		Alpha_3Code:  CountryCodeNIC,
		NumericCode:  558,
		CurrencyCode: pbmoney.CurrencyCodeNIO,
	}
	country166 = &Country{
		Name:         "Netherlands",
		Alpha_2Code:  CountryCodeNL,
		Alpha_3Code:  CountryCodeNLD,
		NumericCode:  528,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country167 = &Country{
		Name:         "Norway",
		Alpha_2Code:  CountryCodeNO,
		Alpha_3Code:  CountryCodeNOR,
		NumericCode:  578,
		CurrencyCode: pbmoney.CurrencyCodeNOK,
	}
	country168 = &Country{
		Name:         "Nepal",
		Alpha_2Code:  CountryCodeNP,
		Alpha_3Code:  CountryCodeNPL,
		NumericCode:  524,
		CurrencyCode: pbmoney.CurrencyCodeNPR,
	}
	country169 = &Country{
		Name:         "Nauru",
		Alpha_2Code:  CountryCodeNR,
		Alpha_3Code:  CountryCodeNRU,
		NumericCode:  520,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country170 = &Country{
		Name:         "Niue",
		Alpha_2Code:  CountryCodeNU,
		Alpha_3Code:  CountryCodeNIU,
		NumericCode:  570,
		CurrencyCode: pbmoney.CurrencyCodeNZD,
	}
	country171 = &Country{
		Name:         "New Zealand",
		Alpha_2Code:  CountryCodeNZ,
		Alpha_3Code:  CountryCodeNZL,
		NumericCode:  554,
		CurrencyCode: pbmoney.CurrencyCodeNZD,
	}
	country172 = &Country{
		Name:         "Oman",
		Alpha_2Code:  CountryCodeOM,
		Alpha_3Code:  CountryCodeOMN,
		NumericCode:  512,
		CurrencyCode: pbmoney.CurrencyCodeOMR,
	}
	country173 = &Country{
		Name:         "Panama",
		Alpha_2Code:  CountryCodePA,
		Alpha_3Code:  CountryCodePAN,
		NumericCode:  591,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country174 = &Country{
		Name:         "Peru",
		Alpha_2Code:  CountryCodePE,
		Alpha_3Code:  CountryCodePER,
		NumericCode:  604,
		CurrencyCode: pbmoney.CurrencyCodePEN,
	}
	country175 = &Country{
		Name:         "French Polynesia",
		Alpha_2Code:  CountryCodePF,
		Alpha_3Code:  CountryCodePYF,
		NumericCode:  258,
		CurrencyCode: pbmoney.CurrencyCodeXPF,
	}
	country176 = &Country{
		Name:         "Papua New Guinea",
		Alpha_2Code:  CountryCodePG,
		Alpha_3Code:  CountryCodePNG,
		NumericCode:  598,
		CurrencyCode: pbmoney.CurrencyCodePGK,
	}
	country177 = &Country{
		Name:         "Philippines",
		Alpha_2Code:  CountryCodePH,
		Alpha_3Code:  CountryCodePHL,
		NumericCode:  608,
		CurrencyCode: pbmoney.CurrencyCodePHP,
	}
	country178 = &Country{
		Name:         "Pakistan",
		Alpha_2Code:  CountryCodePK,
		Alpha_3Code:  CountryCodePAK,
		NumericCode:  586,
		CurrencyCode: pbmoney.CurrencyCodePKR,
	}
	country179 = &Country{
		Name:         "Poland",
		Alpha_2Code:  CountryCodePL,
		Alpha_3Code:  CountryCodePOL,
		NumericCode:  616,
		CurrencyCode: pbmoney.CurrencyCodePLN,
	}
	country180 = &Country{
		Name:         "Saint Pierre and Miquelon",
		Alpha_2Code:  CountryCodePM,
		Alpha_3Code:  CountryCodeSPM,
		NumericCode:  666,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country181 = &Country{
		Name:         "Pitcairn",
		Alpha_2Code:  CountryCodePN,
		Alpha_3Code:  CountryCodePCN,
		NumericCode:  612,
		CurrencyCode: pbmoney.CurrencyCodeNZD,
	}
	country182 = &Country{
		Name:         "Puerto Rico",
		Alpha_2Code:  CountryCodePR,
		Alpha_3Code:  CountryCodePRI,
		NumericCode:  630,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country183 = &Country{
		Name:         "Palestine, State of",
		Alpha_2Code:  CountryCodePS,
		Alpha_3Code:  CountryCodePSE,
		NumericCode:  275,
		CurrencyCode: pbmoney.CurrencyCode_CURRENCY_CODE_NONE,
	}
	country184 = &Country{
		Name:         "Portugal",
		Alpha_2Code:  CountryCodePT,
		Alpha_3Code:  CountryCodePRT,
		NumericCode:  620,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country185 = &Country{
		Name:         "Palau",
		Alpha_2Code:  CountryCodePW,
		Alpha_3Code:  CountryCodePLW,
		NumericCode:  585,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country186 = &Country{
		Name:         "Paraguay",
		Alpha_2Code:  CountryCodePY,
		Alpha_3Code:  CountryCodePRY,
		NumericCode:  600,
		CurrencyCode: pbmoney.CurrencyCodePYG,
	}
	country187 = &Country{
		Name:         "Qatar",
		Alpha_2Code:  CountryCodeQA,
		Alpha_3Code:  CountryCodeQAT,
		NumericCode:  634,
		CurrencyCode: pbmoney.CurrencyCodeQAR,
	}
	country188 = &Country{
		Name:         "Réunion",
		Alpha_2Code:  CountryCodeRE,
		Alpha_3Code:  CountryCodeREU,
		NumericCode:  638,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country189 = &Country{
		Name:         "Romania",
		Alpha_2Code:  CountryCodeRO,
		Alpha_3Code:  CountryCodeROU,
		NumericCode:  642,
		CurrencyCode: pbmoney.CurrencyCodeRON,
	}
	country190 = &Country{
		Name:         "Serbia",
		Alpha_2Code:  CountryCodeRS,
		Alpha_3Code:  CountryCodeSRB,
		NumericCode:  688,
		CurrencyCode: pbmoney.CurrencyCodeRSD,
	}
	country191 = &Country{
		Name:         "Russian Federation",
		Alpha_2Code:  CountryCodeRU,
		Alpha_3Code:  CountryCodeRUS,
		NumericCode:  643,
		CurrencyCode: pbmoney.CurrencyCodeRUB,
	}
	country192 = &Country{
		Name:         "Rwanda",
		Alpha_2Code:  CountryCodeRW,
		Alpha_3Code:  CountryCodeRWA,
		NumericCode:  646,
		CurrencyCode: pbmoney.CurrencyCodeRWF,
	}
	country193 = &Country{
		Name:         "Saudi Arabia",
		Alpha_2Code:  CountryCodeSA,
		Alpha_3Code:  CountryCodeSAU,
		NumericCode:  682,
		CurrencyCode: pbmoney.CurrencyCodeSAR,
	}
	country194 = &Country{
		Name:         "Solomon Islands",
		Alpha_2Code:  CountryCodeSB,
		Alpha_3Code:  CountryCodeSLB,
		NumericCode:  90,
		CurrencyCode: pbmoney.CurrencyCodeSBD,
	}
	country195 = &Country{
		Name:         "Seychelles",
		Alpha_2Code:  CountryCodeSC,
		Alpha_3Code:  CountryCodeSYC,
		NumericCode:  690,
		CurrencyCode: pbmoney.CurrencyCodeSCR,
	}
	country196 = &Country{
		Name:         "Sudan",
		Alpha_2Code:  CountryCodeSD,
		Alpha_3Code:  CountryCodeSDN,
		NumericCode:  729,
		CurrencyCode: pbmoney.CurrencyCodeSDG,
	}
	country197 = &Country{
		Name:         "Sweden",
		Alpha_2Code:  CountryCodeSE,
		Alpha_3Code:  CountryCodeSWE,
		NumericCode:  752,
		CurrencyCode: pbmoney.CurrencyCodeSEK,
	}
	country198 = &Country{
		Name:         "Singapore",
		Alpha_2Code:  CountryCodeSG,
		Alpha_3Code:  CountryCodeSGP,
		NumericCode:  702,
		CurrencyCode: pbmoney.CurrencyCodeSGD,
	}
	country199 = &Country{
		Name:         "Saint Helena, Ascension and Tristan da Cunha",
		Alpha_2Code:  CountryCodeSH,
		Alpha_3Code:  CountryCodeSHN,
		NumericCode:  654,
		CurrencyCode: pbmoney.CurrencyCodeSHP,
	}
	country200 = &Country{
		Name:         "Slovenia",
		Alpha_2Code:  CountryCodeSI,
		Alpha_3Code:  CountryCodeSVN,
		NumericCode:  705,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country201 = &Country{
		Name:         "Svalbard and Jan Mayen",
		Alpha_2Code:  CountryCodeSJ,
		Alpha_3Code:  CountryCodeSJM,
		NumericCode:  744,
		CurrencyCode: pbmoney.CurrencyCodeNOK,
	}
	country202 = &Country{
		Name:         "Slovakia",
		Alpha_2Code:  CountryCodeSK,
		Alpha_3Code:  CountryCodeSVK,
		NumericCode:  703,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country203 = &Country{
		Name:         "Sierra Leone",
		Alpha_2Code:  CountryCodeSL,
		Alpha_3Code:  CountryCodeSLE,
		NumericCode:  694,
		CurrencyCode: pbmoney.CurrencyCodeSLL,
	}
	country204 = &Country{
		Name:         "San Marino",
		Alpha_2Code:  CountryCodeSM,
		Alpha_3Code:  CountryCodeSMR,
		NumericCode:  674,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country205 = &Country{
		Name:         "Senegal",
		Alpha_2Code:  CountryCodeSN,
		Alpha_3Code:  CountryCodeSEN,
		NumericCode:  686,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country206 = &Country{
		Name:         "Somalia",
		Alpha_2Code:  CountryCodeSO,
		Alpha_3Code:  CountryCodeSOM,
		NumericCode:  706,
		CurrencyCode: pbmoney.CurrencyCodeSOS,
	}
	country207 = &Country{
		Name:         "Suriname",
		Alpha_2Code:  CountryCodeSR,
		Alpha_3Code:  CountryCodeSUR,
		NumericCode:  740,
		CurrencyCode: pbmoney.CurrencyCodeSRD,
	}
	country208 = &Country{
		Name:         "South Sudan",
		Alpha_2Code:  CountryCodeSS,
		Alpha_3Code:  CountryCodeSSD,
		NumericCode:  728,
		CurrencyCode: pbmoney.CurrencyCodeSSP,
	}
	country209 = &Country{
		Name:         "Sao Tome and Principe",
		Alpha_2Code:  CountryCodeST,
		Alpha_3Code:  CountryCodeSTP,
		NumericCode:  678,
		CurrencyCode: pbmoney.CurrencyCodeSTD,
	}
	country210 = &Country{
		Name:         "El Salvador",
		Alpha_2Code:  CountryCodeSV,
		Alpha_3Code:  CountryCodeSLV,
		NumericCode:  222,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country211 = &Country{
		Name:         "Sint Maarten (Dutch part)",
		Alpha_2Code:  CountryCodeSX,
		Alpha_3Code:  CountryCodeSXM,
		NumericCode:  534,
		CurrencyCode: pbmoney.CurrencyCodeANG,
	}
	country212 = &Country{
		Name:         "Syrian Arab Republic",
		Alpha_2Code:  CountryCodeSY,
		Alpha_3Code:  CountryCodeSYR,
		NumericCode:  760,
		CurrencyCode: pbmoney.CurrencyCodeSYP,
	}
	country213 = &Country{
		Name:         "Swaziland",
		Alpha_2Code:  CountryCodeSZ,
		Alpha_3Code:  CountryCodeSWZ,
		NumericCode:  748,
		CurrencyCode: pbmoney.CurrencyCodeSZL,
	}
	country214 = &Country{
		Name:         "Turks and Caicos Islands",
		Alpha_2Code:  CountryCodeTC,
		Alpha_3Code:  CountryCodeTCA,
		NumericCode:  796,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country215 = &Country{
		Name:         "Chad",
		Alpha_2Code:  CountryCodeTD,
		Alpha_3Code:  CountryCodeTCD,
		NumericCode:  148,
		CurrencyCode: pbmoney.CurrencyCodeXAF,
	}
	country216 = &Country{
		Name:         "French Southern Territories",
		Alpha_2Code:  CountryCodeTF,
		Alpha_3Code:  CountryCodeATF,
		NumericCode:  260,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country217 = &Country{
		Name:         "Togo",
		Alpha_2Code:  CountryCodeTG,
		Alpha_3Code:  CountryCodeTGO,
		NumericCode:  768,
		CurrencyCode: pbmoney.CurrencyCodeXOF,
	}
	country218 = &Country{
		Name:         "Thailand",
		Alpha_2Code:  CountryCodeTH,
		Alpha_3Code:  CountryCodeTHA,
		NumericCode:  764,
		CurrencyCode: pbmoney.CurrencyCodeTHB,
	}
	country219 = &Country{
		Name:         "Tajikistan",
		Alpha_2Code:  CountryCodeTJ,
		Alpha_3Code:  CountryCodeTJK,
		NumericCode:  762,
		CurrencyCode: pbmoney.CurrencyCodeTJS,
	}
	country220 = &Country{
		Name:         "Tokelau",
		Alpha_2Code:  CountryCodeTK,
		Alpha_3Code:  CountryCodeTKL,
		NumericCode:  772,
		CurrencyCode: pbmoney.CurrencyCodeNZD,
	}
	country221 = &Country{
		Name:         "Timor-Leste",
		Alpha_2Code:  CountryCodeTL,
		Alpha_3Code:  CountryCodeTLS,
		NumericCode:  626,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country222 = &Country{
		Name:         "Turkmenistan",
		Alpha_2Code:  CountryCodeTM,
		Alpha_3Code:  CountryCodeTKM,
		NumericCode:  795,
		CurrencyCode: pbmoney.CurrencyCodeTMT,
	}
	country223 = &Country{
		Name:         "Tunisia",
		Alpha_2Code:  CountryCodeTN,
		Alpha_3Code:  CountryCodeTUN,
		NumericCode:  788,
		CurrencyCode: pbmoney.CurrencyCodeTND,
	}
	country224 = &Country{
		Name:         "Tonga",
		Alpha_2Code:  CountryCodeTO,
		Alpha_3Code:  CountryCodeTON,
		NumericCode:  776,
		CurrencyCode: pbmoney.CurrencyCodeTOP,
	}
	country225 = &Country{
		Name:         "Turkey",
		Alpha_2Code:  CountryCodeTR,
		Alpha_3Code:  CountryCodeTUR,
		NumericCode:  792,
		CurrencyCode: pbmoney.CurrencyCodeTRY,
	}
	country226 = &Country{
		Name:         "Trinidad and Tobago",
		Alpha_2Code:  CountryCodeTT,
		Alpha_3Code:  CountryCodeTTO,
		NumericCode:  780,
		CurrencyCode: pbmoney.CurrencyCodeTTD,
	}
	country227 = &Country{
		Name:         "Tuvalu",
		Alpha_2Code:  CountryCodeTV,
		Alpha_3Code:  CountryCodeTUV,
		NumericCode:  798,
		CurrencyCode: pbmoney.CurrencyCodeAUD,
	}
	country228 = &Country{
		Name:         "Taiwan, Province of China",
		Alpha_2Code:  CountryCodeTW,
		Alpha_3Code:  CountryCodeTWN,
		NumericCode:  158,
		CurrencyCode: pbmoney.CurrencyCodeTWD,
	}
	country229 = &Country{
		Name:         "Tanzania, United Republic of",
		Alpha_2Code:  CountryCodeTZ,
		Alpha_3Code:  CountryCodeTZA,
		NumericCode:  834,
		CurrencyCode: pbmoney.CurrencyCodeTZS,
	}
	country230 = &Country{
		Name:         "Ukraine",
		Alpha_2Code:  CountryCodeUA,
		Alpha_3Code:  CountryCodeUKR,
		NumericCode:  804,
		CurrencyCode: pbmoney.CurrencyCodeUAH,
	}
	country231 = &Country{
		Name:         "Uganda",
		Alpha_2Code:  CountryCodeUG,
		Alpha_3Code:  CountryCodeUGA,
		NumericCode:  800,
		CurrencyCode: pbmoney.CurrencyCodeUGX,
	}
	country232 = &Country{
		Name:         "United States Minor Outlying Islands",
		Alpha_2Code:  CountryCodeUM,
		Alpha_3Code:  CountryCodeUMI,
		NumericCode:  581,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country233 = &Country{
		Name:         "United States",
		Alpha_2Code:  CountryCodeUS,
		Alpha_3Code:  CountryCodeUSA,
		NumericCode:  840,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country234 = &Country{
		Name:         "Uruguay",
		Alpha_2Code:  CountryCodeUY,
		Alpha_3Code:  CountryCodeURY,
		NumericCode:  858,
		CurrencyCode: pbmoney.CurrencyCodeUYU,
	}
	country235 = &Country{
		Name:         "Uzbekistan",
		Alpha_2Code:  CountryCodeUZ,
		Alpha_3Code:  CountryCodeUZB,
		NumericCode:  860,
		CurrencyCode: pbmoney.CurrencyCodeUZS,
	}
	country236 = &Country{
		Name:         "Holy See (Vatican City State)",
		Alpha_2Code:  CountryCodeVA,
		Alpha_3Code:  CountryCodeVAT,
		NumericCode:  336,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country237 = &Country{
		Name:         "Saint Vincent and the Grenadines",
		Alpha_2Code:  CountryCodeVC,
		Alpha_3Code:  CountryCodeVCT,
		NumericCode:  670,
		CurrencyCode: pbmoney.CurrencyCodeXCD,
	}
	country238 = &Country{
		Name:         "Venezuela, Bolivarian Republic of",
		Alpha_2Code:  CountryCodeVE,
		Alpha_3Code:  CountryCodeVEN,
		NumericCode:  862,
		CurrencyCode: pbmoney.CurrencyCodeVEF,
	}
	country239 = &Country{
		Name:         "Virgin Islands, British",
		Alpha_2Code:  CountryCodeVG,
		Alpha_3Code:  CountryCodeVGB,
		NumericCode:  92,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country240 = &Country{
		Name:         "Virgin Islands, U.S.",
		Alpha_2Code:  CountryCodeVI,
		Alpha_3Code:  CountryCodeVIR,
		NumericCode:  850,
		CurrencyCode: pbmoney.CurrencyCodeUSD,
	}
	country241 = &Country{
		Name:         "Viet Nam",
		Alpha_2Code:  CountryCodeVN,
		Alpha_3Code:  CountryCodeVNM,
		NumericCode:  704,
		CurrencyCode: pbmoney.CurrencyCodeVND,
	}
	country242 = &Country{
		Name:         "Vanuatu",
		Alpha_2Code:  CountryCodeVU,
		Alpha_3Code:  CountryCodeVUT,
		NumericCode:  548,
		CurrencyCode: pbmoney.CurrencyCodeVUV,
	}
	country243 = &Country{
		Name:         "Wallis and Futuna",
		Alpha_2Code:  CountryCodeWF,
		Alpha_3Code:  CountryCodeWLF,
		NumericCode:  876,
		CurrencyCode: pbmoney.CurrencyCodeXPF,
	}
	country244 = &Country{
		Name:         "Samoa",
		Alpha_2Code:  CountryCodeWS,
		Alpha_3Code:  CountryCodeWSM,
		NumericCode:  882,
		CurrencyCode: pbmoney.CurrencyCodeWST,
	}
	country245 = &Country{
		Name:         "Yemen",
		Alpha_2Code:  CountryCodeYE,
		Alpha_3Code:  CountryCodeYEM,
		NumericCode:  887,
		CurrencyCode: pbmoney.CurrencyCodeYER,
	}
	country246 = &Country{
		Name:         "Mayotte",
		Alpha_2Code:  CountryCodeYT,
		Alpha_3Code:  CountryCodeMYT,
		NumericCode:  175,
		CurrencyCode: pbmoney.CurrencyCodeEUR,
	}
	country247 = &Country{
		Name:         "South Africa",
		Alpha_2Code:  CountryCodeZA,
		Alpha_3Code:  CountryCodeZAF,
		NumericCode:  710,
		CurrencyCode: pbmoney.CurrencyCodeZAR,
	}
	country248 = &Country{
		Name:         "Zambia",
		Alpha_2Code:  CountryCodeZM,
		Alpha_3Code:  CountryCodeZMB,
		NumericCode:  894,
		CurrencyCode: pbmoney.CurrencyCodeZMW,
	}
	country249 = &Country{
		Name:         "Zimbabwe",
		Alpha_2Code:  CountryCodeZW,
		Alpha_3Code:  CountryCodeZWE,
		NumericCode:  716,
		CurrencyCode: pbmoney.CurrencyCodeZWL,
	}
)

// CountryAlpha2CodeSimpleValueOf returns the value of the simple string s.
func CountryAlpha2CodeSimpleValueOf(s string) (CountryAlpha2Code, error) {
	value, ok := SimpleStringToCountryAlpha2Code[strings.ToUpper(s)]
	if !ok {
		return CountryAlpha2Code_COUNTRY_ALPHA_2_CODE_NONE, fmt.Errorf("pb: no CountryAlpha2Code for %s", s)
	}
	return value, nil
}

// CountryAlpha3CodeSimpleValueOf returns the value of the simple string s.
func CountryAlpha3CodeSimpleValueOf(s string) (CountryAlpha3Code, error) {
	value, ok := SimpleStringToCountryAlpha3Code[strings.ToUpper(s)]
	if !ok {
		return CountryAlpha3Code_COUNTRY_ALPHA_3_CODE_NONE, fmt.Errorf("pb: no CountryAlpha3Code for %s", s)
	}
	return value, nil
}

// SimpleString returns the simple string.
func (c CountryAlpha2Code) SimpleString() string {
	simpleString, ok := CountryAlpha2CodeToSimpleString[c]
	if !ok {
		return ""
	}
	return simpleString
}

// Country returns the associated Country, or nil if no Country is known.
func (c CountryAlpha2Code) Country() *Country {
	country, ok := CountryAlpha2CodeToCountry[c]
	if !ok {
		return nil
	}
	return country
}

// SimpleString returns the simple string.
func (c CountryAlpha3Code) SimpleString() string {
	simpleString, ok := CountryAlpha3CodeToSimpleString[c]
	if !ok {
		return ""
	}
	return simpleString
}

// Country returns the associated Country, or nil if no Country is known.
func (c CountryAlpha3Code) Country() *Country {
	country, ok := CountryAlpha3CodeToCountry[c]
	if !ok {
		return nil
	}
	return country
}

// StreetDirectionSimpleValueOf returns the StreetDirection for the given simple string.
func StreetDirectionSimpleValueOf(s string) (StreetDirection, error) {
	value, ok := StreetDirection_value[fmt.Sprintf("STREET_DIRECTION_%s", strings.ToUpper(s))]
	if !ok {
		return StreetDirection_STREET_DIRECTION_NONE, fmt.Errorf("pb: no StreetDirection for %s", s)
	}
	return StreetDirection(value), nil
}

// SimpleString returns a simple string for the given StreetDirection.
func (s StreetDirection) SimpleString() string {
	value, ok := StreetDirection_name[int32(s)]
	if !ok {
		return ""
	}
	return strings.ToLower(strings.TrimPrefix("STREET_DIRECTION_", value))
}

// SimplePostalAddress returns a SimplePostalAddress for the given PostalAddress.
func (p *PostalAddress) SimplePostalAddress() *SimplePostalAddress {
	buffer := bytes.NewBuffer(nil)
	if p.StreetNumber != 0 {
		_, _ = buffer.WriteString(fmt.Sprintf("%d", p.StreetNumber))
		if p.StreetNumberPostfix != "" {
			_, _ = buffer.WriteString(p.StreetNumberPostfix)
		}
		_ = buffer.WriteByte(' ')
	}
	if p.PreStreetDirection != StreetDirection_STREET_DIRECTION_NONE {
		_, _ = buffer.WriteString(p.PreStreetDirection.SimpleString())
		_ = buffer.WriteByte(' ')
	}
	if p.StreetName != "" {
		_, _ = buffer.WriteString(p.StreetName)
		_ = buffer.WriteByte(' ')
	}
	if p.PostStreetDirection != StreetDirection_STREET_DIRECTION_NONE {
		_, _ = buffer.WriteString(p.PostStreetDirection.SimpleString())
		_ = buffer.WriteByte(' ')
	}
	if p.StreetTypeAbbreviation != "" {
		_, _ = buffer.WriteString(p.StreetTypeAbbreviation)
	}
	streetAddress := buffer.String()

	buffer = bytes.NewBuffer(nil)
	if p.SecondaryAddressTypeAbbreviation != "" {
		_, _ = buffer.WriteString(p.SecondaryAddressTypeAbbreviation)
		_ = buffer.WriteByte(' ')
	}
	if p.SecondaryAddressValue != "" {
		_, _ = buffer.WriteString(p.SecondaryAddressValue)
	}
	streetAddress2 := buffer.String()

	return &SimplePostalAddress{
		StreetAddress:      streetAddress,
		StreetAddress_2:    streetAddress2,
		LocalityName:       p.LocalityName,
		RegionCode:         p.RegionCode,
		PostalCode:         p.PostalCode,
		CountryAlpha_2Code: p.CountryAlpha_2Code,
	}
}
