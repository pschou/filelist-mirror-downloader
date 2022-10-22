// Code generated by "stringer -output types_string.gen.go -type=TagType"; DO NOT EDIT.

package rpm

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RPMTAG_NAME-1000]
	_ = x[RPMTAG_VERSION-1001]
	_ = x[RPMTAG_RELEASE-1002]
	_ = x[RPMTAG_EPOCH-1003]
	_ = x[RPMTAG_SUMMARY-1004]
	_ = x[RPMTAG_DESCRIPTION-1005]
	_ = x[RPMTAG_BUILDTIME-1006]
	_ = x[RPMTAG_BUILDHOST-1007]
	_ = x[RPMTAG_INSTALLTIME-1008]
	_ = x[RPMTAG_SIZE-1009]
	_ = x[RPMTAG_DISTRIBUTION-1010]
	_ = x[RPMTAG_VENDOR-1011]
	_ = x[RPMTAG_GIF-1012]
	_ = x[RPMTAG_XPM-1013]
	_ = x[RPMTAG_LICENSE-1014]
	_ = x[RPMTAG_PACKAGER-1015]
	_ = x[RPMTAG_GROUP-1016]
	_ = x[RPMTAG_CHANGELOG-1017]
	_ = x[RPMTAG_SOURCE-1018]
	_ = x[RPMTAG_PATCH-1019]
	_ = x[RPMTAG_URL-1020]
	_ = x[RPMTAG_OS-1021]
	_ = x[RPMTAG_ARCH-1022]
	_ = x[RPMTAG_PREIN-1023]
	_ = x[RPMTAG_POSTIN-1024]
	_ = x[RPMTAG_PREUN-1025]
	_ = x[RPMTAG_POSTUN-1026]
	_ = x[RPMTAG_OLDFILENAMES-1027]
	_ = x[RPMTAG_FILESIZES-1028]
	_ = x[RPMTAG_FILESTATES-1029]
	_ = x[RPMTAG_FILEMODES-1030]
	_ = x[RPMTAG_FILEUIDS-1031]
	_ = x[RPMTAG_FILEGIDS-1032]
	_ = x[RPMTAG_FILERDEVS-1033]
	_ = x[RPMTAG_FILEMTIMES-1034]
	_ = x[RPMTAG_FILEDIGESTS-1035]
	_ = x[RPMTAG_FILELINKTOS-1036]
	_ = x[RPMTAG_FILEFLAGS-1037]
	_ = x[RPMTAG_ROOT-1038]
	_ = x[RPMTAG_FILEUSERNAME-1039]
	_ = x[RPMTAG_FILEGROUPNAME-1040]
	_ = x[RPMTAG_EXCLUDE-1041]
	_ = x[RPMTAG_EXCLUSIVE-1042]
	_ = x[RPMTAG_ICON-1043]
	_ = x[RPMTAG_SOURCERPM-1044]
	_ = x[RPMTAG_FILEVERIFYFLAGS-1045]
	_ = x[RPMTAG_ARCHIVESIZE-1046]
	_ = x[RPMTAG_PROVIDENAME-1047]
	_ = x[RPMTAG_REQUIREFLAGS-1048]
	_ = x[RPMTAG_REQUIRENAME-1049]
	_ = x[RPMTAG_REQUIREVERSION-1050]
	_ = x[RPMTAG_NOSOURCE-1051]
	_ = x[RPMTAG_NOPATCH-1052]
	_ = x[RPMTAG_CONFLICTFLAGS-1053]
	_ = x[RPMTAG_CONFLICTNAME-1054]
	_ = x[RPMTAG_CONFLICTVERSION-1055]
	_ = x[RPMTAG_DEFAULTPREFIX-1056]
	_ = x[RPMTAG_BUILDROOT-1057]
	_ = x[RPMTAG_INSTALLPREFIX-1058]
	_ = x[RPMTAG_EXCLUDEARCH-1059]
	_ = x[RPMTAG_EXCLUDEOS-1060]
	_ = x[RPMTAG_EXCLUSIVEARCH-1061]
	_ = x[RPMTAG_EXCLUSIVEOS-1062]
	_ = x[RPMTAG_AUTOREQPROV-1063]
	_ = x[RPMTAG_RPMVERSION-1064]
	_ = x[RPMTAG_TRIGGERSCRIPTS-1065]
	_ = x[RPMTAG_TRIGGERNAME-1066]
	_ = x[RPMTAG_TRIGGERVERSION-1067]
	_ = x[RPMTAG_TRIGGERFLAGS-1068]
	_ = x[RPMTAG_TRIGGERINDEX-1069]
	_ = x[RPMTAG_VERIFYSCRIPT-1079]
	_ = x[RPMTAG_CHANGELOGTIME-1080]
	_ = x[RPMTAG_CHANGELOGNAME-1081]
	_ = x[RPMTAG_CHANGELOGTEXT-1082]
	_ = x[RPMTAG_BROKENMD5-1083]
	_ = x[RPMTAG_PREREQ-1084]
	_ = x[RPMTAG_PREINPROG-1085]
	_ = x[RPMTAG_POSTINPROG-1086]
	_ = x[RPMTAG_PREUNPROG-1087]
	_ = x[RPMTAG_POSTUNPROG-1088]
	_ = x[RPMTAG_BUILDARCHS-1089]
	_ = x[RPMTAG_OBSOLETENAME-1090]
	_ = x[RPMTAG_VERIFYSCRIPTPROG-1091]
	_ = x[RPMTAG_TRIGGERSCRIPTPROG-1092]
	_ = x[RPMTAG_DOCDIR-1093]
	_ = x[RPMTAG_COOKIE-1094]
	_ = x[RPMTAG_FILEDEVICES-1095]
	_ = x[RPMTAG_FILEINODES-1096]
	_ = x[RPMTAG_FILELANGS-1097]
	_ = x[RPMTAG_PREFIXES-1098]
	_ = x[RPMTAG_INSTPREFIXES-1099]
	_ = x[RPMTAG_TRIGGERIN-1100]
	_ = x[RPMTAG_TRIGGERUN-1101]
	_ = x[RPMTAG_TRIGGERPOSTUN-1102]
	_ = x[RPMTAG_AUTOREQ-1103]
	_ = x[RPMTAG_AUTOPROV-1104]
	_ = x[RPMTAG_CAPABILITY-1105]
	_ = x[RPMTAG_SOURCEPACKAGE-1106]
	_ = x[RPMTAG_OLDORIGFILENAMES-1107]
	_ = x[RPMTAG_BUILDPREREQ-1108]
	_ = x[RPMTAG_BUILDREQUIRES-1109]
	_ = x[RPMTAG_BUILDCONFLICTS-1110]
	_ = x[RPMTAG_BUILDMACROS-1111]
	_ = x[RPMTAG_PROVIDEFLAGS-1112]
	_ = x[RPMTAG_PROVIDEVERSION-1113]
	_ = x[RPMTAG_OBSOLETEFLAGS-1114]
	_ = x[RPMTAG_OBSOLETEVERSION-1115]
	_ = x[RPMTAG_DIRINDEXES-1116]
	_ = x[RPMTAG_BASENAMES-1117]
	_ = x[RPMTAG_DIRNAMES-1118]
	_ = x[RPMTAG_ORIGDIRINDEXES-1119]
	_ = x[RPMTAG_ORIGBASENAMES-1120]
	_ = x[RPMTAG_ORIGDIRNAMES-1121]
	_ = x[RPMTAG_OPTFLAGS-1122]
	_ = x[RPMTAG_DISTURL-1123]
	_ = x[RPMTAG_PAYLOADFORMAT-1124]
	_ = x[RPMTAG_PAYLOADCOMPRESSOR-1125]
	_ = x[RPMTAG_PAYLOADFLAGS-1126]
	_ = x[RPMTAG_INSTALLCOLOR-1127]
	_ = x[RPMTAG_INSTALLTID-1128]
	_ = x[RPMTAG_REMOVETID-1129]
	_ = x[RPMTAG_SHA1RHN-1130]
	_ = x[RPMTAG_RHNPLATFORM-1131]
	_ = x[RPMTAG_PLATFORM-1132]
	_ = x[RPMTAG_PATCHESNAME-1133]
	_ = x[RPMTAG_PATCHESFLAGS-1134]
	_ = x[RPMTAG_PATCHESVERSION-1135]
	_ = x[RPMTAG_CACHECTIME-1136]
	_ = x[RPMTAG_CACHEPKGPATH-1137]
	_ = x[RPMTAG_CACHEPKGSIZE-1138]
	_ = x[RPMTAG_CACHEPKGMTIME-1139]
	_ = x[RPMTAG_FILECOLORS-1140]
	_ = x[RPMTAG_FILECLASS-1141]
	_ = x[RPMTAG_CLASSDICT-1142]
	_ = x[RPMTAG_FILEDEPENDSX-1143]
	_ = x[RPMTAG_FILEDEPENDSN-1144]
	_ = x[RPMTAG_DEPENDSDICT-1145]
	_ = x[RPMTAG_SOURCEPKGID-1146]
	_ = x[RPMTAG_FILECONTEXTS-1147]
	_ = x[RPMTAG_FSCONTEXTS-1148]
	_ = x[RPMTAG_RECONTEXTS-1149]
	_ = x[RPMTAG_POLICIES-1150]
	_ = x[RPMTAG_PRETRANS-1151]
	_ = x[RPMTAG_POSTTRANS-1152]
	_ = x[RPMTAG_PRETRANSPROG-1153]
	_ = x[RPMTAG_POSTTRANSPROG-1154]
	_ = x[RPMTAG_DISTTAG-1155]
	_ = x[RPMTAG_OLDSUGGESTSNAME-1156]
	_ = x[RPMTAG_OLDSUGGESTSVERSION-1157]
	_ = x[RPMTAG_OLDSUGGESTSFLAGS-1158]
	_ = x[RPMTAG_OLDENHANCESNAME-1159]
	_ = x[RPMTAG_OLDENHANCESVERSION-1160]
	_ = x[RPMTAG_OLDENHANCESFLAGS-1161]
	_ = x[RPMTAG_PRIORITY-1162]
	_ = x[RPMTAG_CVSID-1163]
	_ = x[RPMTAG_BLINKPKGID-1164]
	_ = x[RPMTAG_BLINKHDRID-1165]
	_ = x[RPMTAG_BLINKNEVRA-1166]
	_ = x[RPMTAG_FLINKPKGID-1167]
	_ = x[RPMTAG_FLINKHDRID-1168]
	_ = x[RPMTAG_FLINKNEVRA-1169]
	_ = x[RPMTAG_PACKAGEORIGIN-1170]
	_ = x[RPMTAG_TRIGGERPREIN-1171]
	_ = x[RPMTAG_BUILDSUGGESTS-1172]
	_ = x[RPMTAG_BUILDENHANCES-1173]
	_ = x[RPMTAG_SCRIPTSTATES-1174]
	_ = x[RPMTAG_SCRIPTMETRICS-1175]
	_ = x[RPMTAG_BUILDCPUCLOCK-1176]
	_ = x[RPMTAG_FILEDIGESTALGOS-1177]
	_ = x[RPMTAG_VARIANTS-1178]
	_ = x[RPMTAG_XMAJOR-1179]
	_ = x[RPMTAG_XMINOR-1180]
	_ = x[RPMTAG_REPOTAG-1181]
	_ = x[RPMTAG_KEYWORDS-1182]
	_ = x[RPMTAG_BUILDPLATFORMS-1183]
	_ = x[RPMTAG_PACKAGECOLOR-1184]
	_ = x[RPMTAG_PACKAGEPREFCOLOR-1185]
	_ = x[RPMTAG_XATTRSDICT-1186]
	_ = x[RPMTAG_FILEXATTRSX-1187]
	_ = x[RPMTAG_DEPATTRSDICT-1188]
	_ = x[RPMTAG_CONFLICTATTRSX-1189]
	_ = x[RPMTAG_OBSOLETEATTRSX-1190]
	_ = x[RPMTAG_PROVIDEATTRSX-1191]
	_ = x[RPMTAG_REQUIREATTRSX-1192]
	_ = x[RPMTAG_BUILDPROVIDES-1193]
	_ = x[RPMTAG_BUILDOBSOLETES-1194]
	_ = x[RPMTAG_DBINSTANCE-1195]
	_ = x[RPMTAG_NVRA-1196]
	_ = x[RPMTAG_FILENAMES-5000]
	_ = x[RPMTAG_FILEPROVIDE-5001]
	_ = x[RPMTAG_FILEREQUIRE-5002]
	_ = x[RPMTAG_FSNAMES-5003]
	_ = x[RPMTAG_FSSIZES-5004]
	_ = x[RPMTAG_TRIGGERCONDS-5005]
	_ = x[RPMTAG_TRIGGERTYPE-5006]
	_ = x[RPMTAG_ORIGFILENAMES-5007]
	_ = x[RPMTAG_LONGFILESIZES-5008]
	_ = x[RPMTAG_LONGSIZE-5009]
	_ = x[RPMTAG_FILECAPS-5010]
	_ = x[RPMTAG_FILEDIGESTALGO-5011]
	_ = x[RPMTAG_BUGURL-5012]
	_ = x[RPMTAG_EVR-5013]
	_ = x[RPMTAG_NVR-5014]
	_ = x[RPMTAG_NEVR-5015]
	_ = x[RPMTAG_NEVRA-5016]
	_ = x[RPMTAG_HEADERCOLOR-5017]
	_ = x[RPMTAG_VERBOSE-5018]
	_ = x[RPMTAG_EPOCHNUM-5019]
	_ = x[RPMTAG_PREINFLAGS-5020]
	_ = x[RPMTAG_POSTINFLAGS-5021]
	_ = x[RPMTAG_PREUNFLAGS-5022]
	_ = x[RPMTAG_POSTUNFLAGS-5023]
	_ = x[RPMTAG_PRETRANSFLAGS-5024]
	_ = x[RPMTAG_POSTTRANSFLAGS-5025]
	_ = x[RPMTAG_VERIFYSCRIPTFLAGS-5026]
	_ = x[RPMTAG_TRIGGERSCRIPTFLAGS-5027]
	_ = x[RPMTAG_COLLECTIONS-5029]
	_ = x[RPMTAG_POLICYNAMES-5030]
	_ = x[RPMTAG_POLICYTYPES-5031]
	_ = x[RPMTAG_POLICYTYPESINDEXES-5032]
	_ = x[RPMTAG_POLICYFLAGS-5033]
	_ = x[RPMTAG_VCS-5034]
	_ = x[RPMTAG_ORDERNAME-5035]
	_ = x[RPMTAG_ORDERVERSION-5036]
	_ = x[RPMTAG_ORDERFLAGS-5037]
	_ = x[RPMTAG_MSSFMANIFEST-5038]
	_ = x[RPMTAG_MSSFDOMAIN-5039]
	_ = x[RPMTAG_INSTFILENAMES-5040]
	_ = x[RPMTAG_REQUIRENEVRS-5041]
	_ = x[RPMTAG_PROVIDENEVRS-5042]
	_ = x[RPMTAG_OBSOLETENEVRS-5043]
	_ = x[RPMTAG_CONFLICTNEVRS-5044]
	_ = x[RPMTAG_FILENLINKS-5045]
	_ = x[RPMTAG_RECOMMENDNAME-5046]
	_ = x[RPMTAG_RECOMMENDVERSION-5047]
	_ = x[RPMTAG_RECOMMENDFLAGS-5048]
	_ = x[RPMTAG_SUGGESTNAME-5049]
	_ = x[RPMTAG_SUGGESTVERSION-5050]
	_ = x[RPMTAG_SUGGESTFLAGS-5051]
	_ = x[RPMTAG_SUPPLEMENTNAME-5052]
	_ = x[RPMTAG_SUPPLEMENTVERSION-5053]
	_ = x[RPMTAG_SUPPLEMENTFLAGS-5054]
	_ = x[RPMTAG_ENHANCENAME-5055]
	_ = x[RPMTAG_ENHANCEVERSION-5056]
	_ = x[RPMTAG_ENHANCEFLAGS-5057]
	_ = x[RPMTAG_RECOMMENDNEVRS-5058]
	_ = x[RPMTAG_SUGGESTNEVRS-5059]
	_ = x[RPMTAG_SUPPLEMENTNEVRS-5060]
	_ = x[RPMTAG_ENHANCENEVRS-5061]
	_ = x[RPMTAG_ENCODING-5062]
	_ = x[RPMTAG_FILETRIGGERIN-5063]
	_ = x[RPMTAG_FILETRIGGERUN-5064]
	_ = x[RPMTAG_FILETRIGGERPOSTUN-5065]
	_ = x[RPMTAG_FILETRIGGERSCRIPTS-5066]
	_ = x[RPMTAG_FILETRIGGERSCRIPTPROG-5067]
	_ = x[RPMTAG_FILETRIGGERSCRIPTFLAGS-5068]
	_ = x[RPMTAG_FILETRIGGERNAME-5069]
	_ = x[RPMTAG_FILETRIGGERINDEX-5070]
	_ = x[RPMTAG_FILETRIGGERVERSION-5071]
	_ = x[RPMTAG_FILETRIGGERFLAGS-5072]
	_ = x[RPMTAG_TRANSFILETRIGGERIN-5073]
	_ = x[RPMTAG_TRANSFILETRIGGERUN-5074]
	_ = x[RPMTAG_TRANSFILETRIGGERPOSTUN-5075]
	_ = x[RPMTAG_TRANSFILETRIGGERSCRIPTS-5076]
	_ = x[RPMTAG_TRANSFILETRIGGERSCRIPTPROG-5077]
	_ = x[RPMTAG_TRANSFILETRIGGERSCRIPTFLAGS-5078]
	_ = x[RPMTAG_TRANSFILETRIGGERNAME-5079]
	_ = x[RPMTAG_TRANSFILETRIGGERINDEX-5080]
	_ = x[RPMTAG_TRANSFILETRIGGERVERSION-5081]
	_ = x[RPMTAG_TRANSFILETRIGGERFLAGS-5082]
	_ = x[RPMTAG_REMOVEPATHPOSTFIXES-5083]
	_ = x[RPMTAG_FILETRIGGERPRIORITIES-5084]
	_ = x[RPMTAG_TRANSFILETRIGGERPRIORITIES-5085]
	_ = x[RPMTAG_FILETRIGGERCONDS-5086]
	_ = x[RPMTAG_FILETRIGGERTYPE-5087]
	_ = x[RPMTAG_TRANSFILETRIGGERCONDS-5088]
	_ = x[RPMTAG_TRANSFILETRIGGERTYPE-5089]
	_ = x[RPMTAG_FILESIGNATURES-5090]
	_ = x[RPMTAG_FILESIGNATURELENGTH-5091]
	_ = x[RPMTAG_PAYLOADDIGEST-5092]
	_ = x[RPMTAG_PAYLOADDIGESTALGO-5093]
	_ = x[RPMTAG_AUTOINSTALLED-5094]
	_ = x[RPMTAG_IDENTITY-5095]
	_ = x[RPMTAG_MODULARITYLABEL-5096]
	_ = x[RPMTAG_PAYLOADDIGESTALT-5097]
	_ = x[RPMTAG_HEADERI18NTABLE-100]
	_ = x[RPMTAG_HEADERIMAGE-61]
	_ = x[RPMTAG_HEADERIMMUTABLE-63]
	_ = x[RPMTAG_HEADERREGIONS-64]
	_ = x[RPMTAG_SIG_BASE-256]
	_ = x[RPMTAG_HEADERSIGNATURES-62]
	_ = x[RPMTAG_SIGSIZE-257]
	_ = x[RPMTAG_PUBKEYS-266]
	_ = x[RPMTAG_DSAHEADER-267]
	_ = x[RPMTAG_RSAHEADER-268]
	_ = x[RPMTAG_SHA1HEADER-269]
	_ = x[RPMTAG_LONGSIGSIZE-270]
	_ = x[RPMTAG_LONGARCHIVESIZE-271]
	_ = x[RPMTAG_SHA256HEADER-273]
	_ = x[RPMTAG_SIGLEMD5_1-258]
	_ = x[RPMTAG_SIGPGP-259]
	_ = x[RPMTAG_SIGLEMD5_2-260]
	_ = x[RPMTAG_SIGMD5-261]
	_ = x[RPMTAG_SIGGPG-262]
	_ = x[RPMTAG_SIGPGP5-263]
	_ = x[RPMTAG_BADSHA1_1-264]
	_ = x[RPMTAG_BADSHA1_2-265]
	_ = x[RPMTAG_C-1054]
	_ = x[RPMTAG_SVNID-1163]
	_ = x[RPMTAG_ENHANCES-5055]
	_ = x[RPMTAG_E-1003]
	_ = x[RPMTAG_FILEMD5S-1035]
	_ = x[RPMTAG_N-1000]
	_ = x[RPMTAG_O-1090]
	_ = x[RPMTAG_OLDENHANCES-1159]
	_ = x[RPMTAG_OLDSUGGESTS-1156]
	_ = x[RPMTAG_P-1047]
	_ = x[RPMTAG_RECOMMENDS-5046]
	_ = x[RPMTAG_R-1002]
	_ = x[RPMTAG_REQUIRES-1049]
	_ = x[RPMTAG_HDRID-269]
	_ = x[RPMTAG_PKGID-261]
	_ = x[RPMTAG_SUGGESTS-5049]
	_ = x[RPMTAG_SUPPLEMENTS-5052]
	_ = x[RPMTAG_V-1001]
}

const (
	_TagType_name_0 = "RPMTAG_HEADERIMAGERPMTAG_HEADERSIGNATURESRPMTAG_HEADERIMMUTABLERPMTAG_HEADERREGIONS"
	_TagType_name_1 = "RPMTAG_HEADERI18NTABLE"
	_TagType_name_2 = "RPMTAG_SIG_BASERPMTAG_SIGSIZERPMTAG_SIGLEMD5_1RPMTAG_SIGPGPRPMTAG_SIGLEMD5_2RPMTAG_SIGMD5RPMTAG_SIGGPGRPMTAG_SIGPGP5RPMTAG_BADSHA1_1RPMTAG_BADSHA1_2RPMTAG_PUBKEYSRPMTAG_DSAHEADERRPMTAG_RSAHEADERRPMTAG_SHA1HEADERRPMTAG_LONGSIGSIZERPMTAG_LONGARCHIVESIZE"
	_TagType_name_3 = "RPMTAG_SHA256HEADER"
	_TagType_name_4 = "RPMTAG_NAMERPMTAG_VERSIONRPMTAG_RELEASERPMTAG_EPOCHRPMTAG_SUMMARYRPMTAG_DESCRIPTIONRPMTAG_BUILDTIMERPMTAG_BUILDHOSTRPMTAG_INSTALLTIMERPMTAG_SIZERPMTAG_DISTRIBUTIONRPMTAG_VENDORRPMTAG_GIFRPMTAG_XPMRPMTAG_LICENSERPMTAG_PACKAGERRPMTAG_GROUPRPMTAG_CHANGELOGRPMTAG_SOURCERPMTAG_PATCHRPMTAG_URLRPMTAG_OSRPMTAG_ARCHRPMTAG_PREINRPMTAG_POSTINRPMTAG_PREUNRPMTAG_POSTUNRPMTAG_OLDFILENAMESRPMTAG_FILESIZESRPMTAG_FILESTATESRPMTAG_FILEMODESRPMTAG_FILEUIDSRPMTAG_FILEGIDSRPMTAG_FILERDEVSRPMTAG_FILEMTIMESRPMTAG_FILEDIGESTSRPMTAG_FILELINKTOSRPMTAG_FILEFLAGSRPMTAG_ROOTRPMTAG_FILEUSERNAMERPMTAG_FILEGROUPNAMERPMTAG_EXCLUDERPMTAG_EXCLUSIVERPMTAG_ICONRPMTAG_SOURCERPMRPMTAG_FILEVERIFYFLAGSRPMTAG_ARCHIVESIZERPMTAG_PROVIDENAMERPMTAG_REQUIREFLAGSRPMTAG_REQUIRENAMERPMTAG_REQUIREVERSIONRPMTAG_NOSOURCERPMTAG_NOPATCHRPMTAG_CONFLICTFLAGSRPMTAG_CONFLICTNAMERPMTAG_CONFLICTVERSIONRPMTAG_DEFAULTPREFIXRPMTAG_BUILDROOTRPMTAG_INSTALLPREFIXRPMTAG_EXCLUDEARCHRPMTAG_EXCLUDEOSRPMTAG_EXCLUSIVEARCHRPMTAG_EXCLUSIVEOSRPMTAG_AUTOREQPROVRPMTAG_RPMVERSIONRPMTAG_TRIGGERSCRIPTSRPMTAG_TRIGGERNAMERPMTAG_TRIGGERVERSIONRPMTAG_TRIGGERFLAGSRPMTAG_TRIGGERINDEX"
	_TagType_name_5 = "RPMTAG_VERIFYSCRIPTRPMTAG_CHANGELOGTIMERPMTAG_CHANGELOGNAMERPMTAG_CHANGELOGTEXTRPMTAG_BROKENMD5RPMTAG_PREREQRPMTAG_PREINPROGRPMTAG_POSTINPROGRPMTAG_PREUNPROGRPMTAG_POSTUNPROGRPMTAG_BUILDARCHSRPMTAG_OBSOLETENAMERPMTAG_VERIFYSCRIPTPROGRPMTAG_TRIGGERSCRIPTPROGRPMTAG_DOCDIRRPMTAG_COOKIERPMTAG_FILEDEVICESRPMTAG_FILEINODESRPMTAG_FILELANGSRPMTAG_PREFIXESRPMTAG_INSTPREFIXESRPMTAG_TRIGGERINRPMTAG_TRIGGERUNRPMTAG_TRIGGERPOSTUNRPMTAG_AUTOREQRPMTAG_AUTOPROVRPMTAG_CAPABILITYRPMTAG_SOURCEPACKAGERPMTAG_OLDORIGFILENAMESRPMTAG_BUILDPREREQRPMTAG_BUILDREQUIRESRPMTAG_BUILDCONFLICTSRPMTAG_BUILDMACROSRPMTAG_PROVIDEFLAGSRPMTAG_PROVIDEVERSIONRPMTAG_OBSOLETEFLAGSRPMTAG_OBSOLETEVERSIONRPMTAG_DIRINDEXESRPMTAG_BASENAMESRPMTAG_DIRNAMESRPMTAG_ORIGDIRINDEXESRPMTAG_ORIGBASENAMESRPMTAG_ORIGDIRNAMESRPMTAG_OPTFLAGSRPMTAG_DISTURLRPMTAG_PAYLOADFORMATRPMTAG_PAYLOADCOMPRESSORRPMTAG_PAYLOADFLAGSRPMTAG_INSTALLCOLORRPMTAG_INSTALLTIDRPMTAG_REMOVETIDRPMTAG_SHA1RHNRPMTAG_RHNPLATFORMRPMTAG_PLATFORMRPMTAG_PATCHESNAMERPMTAG_PATCHESFLAGSRPMTAG_PATCHESVERSIONRPMTAG_CACHECTIMERPMTAG_CACHEPKGPATHRPMTAG_CACHEPKGSIZERPMTAG_CACHEPKGMTIMERPMTAG_FILECOLORSRPMTAG_FILECLASSRPMTAG_CLASSDICTRPMTAG_FILEDEPENDSXRPMTAG_FILEDEPENDSNRPMTAG_DEPENDSDICTRPMTAG_SOURCEPKGIDRPMTAG_FILECONTEXTSRPMTAG_FSCONTEXTSRPMTAG_RECONTEXTSRPMTAG_POLICIESRPMTAG_PRETRANSRPMTAG_POSTTRANSRPMTAG_PRETRANSPROGRPMTAG_POSTTRANSPROGRPMTAG_DISTTAGRPMTAG_OLDSUGGESTSNAMERPMTAG_OLDSUGGESTSVERSIONRPMTAG_OLDSUGGESTSFLAGSRPMTAG_OLDENHANCESNAMERPMTAG_OLDENHANCESVERSIONRPMTAG_OLDENHANCESFLAGSRPMTAG_PRIORITYRPMTAG_CVSIDRPMTAG_BLINKPKGIDRPMTAG_BLINKHDRIDRPMTAG_BLINKNEVRARPMTAG_FLINKPKGIDRPMTAG_FLINKHDRIDRPMTAG_FLINKNEVRARPMTAG_PACKAGEORIGINRPMTAG_TRIGGERPREINRPMTAG_BUILDSUGGESTSRPMTAG_BUILDENHANCESRPMTAG_SCRIPTSTATESRPMTAG_SCRIPTMETRICSRPMTAG_BUILDCPUCLOCKRPMTAG_FILEDIGESTALGOSRPMTAG_VARIANTSRPMTAG_XMAJORRPMTAG_XMINORRPMTAG_REPOTAGRPMTAG_KEYWORDSRPMTAG_BUILDPLATFORMSRPMTAG_PACKAGECOLORRPMTAG_PACKAGEPREFCOLORRPMTAG_XATTRSDICTRPMTAG_FILEXATTRSXRPMTAG_DEPATTRSDICTRPMTAG_CONFLICTATTRSXRPMTAG_OBSOLETEATTRSXRPMTAG_PROVIDEATTRSXRPMTAG_REQUIREATTRSXRPMTAG_BUILDPROVIDESRPMTAG_BUILDOBSOLETESRPMTAG_DBINSTANCERPMTAG_NVRA"
	_TagType_name_6 = "RPMTAG_FILENAMESRPMTAG_FILEPROVIDERPMTAG_FILEREQUIRERPMTAG_FSNAMESRPMTAG_FSSIZESRPMTAG_TRIGGERCONDSRPMTAG_TRIGGERTYPERPMTAG_ORIGFILENAMESRPMTAG_LONGFILESIZESRPMTAG_LONGSIZERPMTAG_FILECAPSRPMTAG_FILEDIGESTALGORPMTAG_BUGURLRPMTAG_EVRRPMTAG_NVRRPMTAG_NEVRRPMTAG_NEVRARPMTAG_HEADERCOLORRPMTAG_VERBOSERPMTAG_EPOCHNUMRPMTAG_PREINFLAGSRPMTAG_POSTINFLAGSRPMTAG_PREUNFLAGSRPMTAG_POSTUNFLAGSRPMTAG_PRETRANSFLAGSRPMTAG_POSTTRANSFLAGSRPMTAG_VERIFYSCRIPTFLAGSRPMTAG_TRIGGERSCRIPTFLAGS"
	_TagType_name_7 = "RPMTAG_COLLECTIONSRPMTAG_POLICYNAMESRPMTAG_POLICYTYPESRPMTAG_POLICYTYPESINDEXESRPMTAG_POLICYFLAGSRPMTAG_VCSRPMTAG_ORDERNAMERPMTAG_ORDERVERSIONRPMTAG_ORDERFLAGSRPMTAG_MSSFMANIFESTRPMTAG_MSSFDOMAINRPMTAG_INSTFILENAMESRPMTAG_REQUIRENEVRSRPMTAG_PROVIDENEVRSRPMTAG_OBSOLETENEVRSRPMTAG_CONFLICTNEVRSRPMTAG_FILENLINKSRPMTAG_RECOMMENDNAMERPMTAG_RECOMMENDVERSIONRPMTAG_RECOMMENDFLAGSRPMTAG_SUGGESTNAMERPMTAG_SUGGESTVERSIONRPMTAG_SUGGESTFLAGSRPMTAG_SUPPLEMENTNAMERPMTAG_SUPPLEMENTVERSIONRPMTAG_SUPPLEMENTFLAGSRPMTAG_ENHANCENAMERPMTAG_ENHANCEVERSIONRPMTAG_ENHANCEFLAGSRPMTAG_RECOMMENDNEVRSRPMTAG_SUGGESTNEVRSRPMTAG_SUPPLEMENTNEVRSRPMTAG_ENHANCENEVRSRPMTAG_ENCODINGRPMTAG_FILETRIGGERINRPMTAG_FILETRIGGERUNRPMTAG_FILETRIGGERPOSTUNRPMTAG_FILETRIGGERSCRIPTSRPMTAG_FILETRIGGERSCRIPTPROGRPMTAG_FILETRIGGERSCRIPTFLAGSRPMTAG_FILETRIGGERNAMERPMTAG_FILETRIGGERINDEXRPMTAG_FILETRIGGERVERSIONRPMTAG_FILETRIGGERFLAGSRPMTAG_TRANSFILETRIGGERINRPMTAG_TRANSFILETRIGGERUNRPMTAG_TRANSFILETRIGGERPOSTUNRPMTAG_TRANSFILETRIGGERSCRIPTSRPMTAG_TRANSFILETRIGGERSCRIPTPROGRPMTAG_TRANSFILETRIGGERSCRIPTFLAGSRPMTAG_TRANSFILETRIGGERNAMERPMTAG_TRANSFILETRIGGERINDEXRPMTAG_TRANSFILETRIGGERVERSIONRPMTAG_TRANSFILETRIGGERFLAGSRPMTAG_REMOVEPATHPOSTFIXESRPMTAG_FILETRIGGERPRIORITIESRPMTAG_TRANSFILETRIGGERPRIORITIESRPMTAG_FILETRIGGERCONDSRPMTAG_FILETRIGGERTYPERPMTAG_TRANSFILETRIGGERCONDSRPMTAG_TRANSFILETRIGGERTYPERPMTAG_FILESIGNATURESRPMTAG_FILESIGNATURELENGTHRPMTAG_PAYLOADDIGESTRPMTAG_PAYLOADDIGESTALGORPMTAG_AUTOINSTALLEDRPMTAG_IDENTITYRPMTAG_MODULARITYLABELRPMTAG_PAYLOADDIGESTALT"
)

var (
	_TagType_index_0 = [...]uint8{0, 18, 41, 63, 83}
	_TagType_index_2 = [...]uint8{0, 15, 29, 46, 59, 76, 89, 102, 116, 132, 148, 162, 178, 194, 211, 229, 251}
	_TagType_index_4 = [...]uint16{0, 11, 25, 39, 51, 65, 83, 99, 115, 133, 144, 163, 176, 186, 196, 210, 225, 237, 253, 266, 278, 288, 297, 308, 320, 333, 345, 358, 377, 393, 410, 426, 441, 456, 472, 489, 507, 525, 541, 552, 571, 591, 605, 621, 632, 648, 670, 688, 706, 725, 743, 764, 779, 793, 813, 832, 854, 874, 890, 910, 928, 944, 964, 982, 1000, 1017, 1038, 1056, 1077, 1096, 1115}
	_TagType_index_5 = [...]uint16{0, 19, 39, 59, 79, 95, 108, 124, 141, 157, 174, 191, 210, 233, 257, 270, 283, 301, 318, 334, 349, 368, 384, 400, 420, 434, 449, 466, 486, 509, 527, 547, 568, 586, 605, 626, 646, 668, 685, 701, 716, 737, 757, 776, 791, 805, 825, 849, 868, 887, 904, 920, 934, 952, 967, 985, 1004, 1025, 1042, 1061, 1080, 1100, 1117, 1133, 1149, 1168, 1187, 1205, 1223, 1242, 1259, 1276, 1291, 1306, 1322, 1341, 1361, 1375, 1397, 1422, 1445, 1467, 1492, 1515, 1530, 1542, 1559, 1576, 1593, 1610, 1627, 1644, 1664, 1683, 1703, 1723, 1742, 1762, 1782, 1804, 1819, 1832, 1845, 1859, 1874, 1895, 1914, 1937, 1954, 1972, 1991, 2012, 2033, 2053, 2073, 2093, 2114, 2131, 2142}
	_TagType_index_6 = [...]uint16{0, 16, 34, 52, 66, 80, 99, 117, 137, 157, 172, 187, 208, 221, 231, 241, 252, 264, 282, 296, 311, 328, 346, 363, 381, 401, 422, 446, 471}
	_TagType_index_7 = [...]uint16{0, 18, 36, 54, 79, 97, 107, 123, 142, 159, 178, 195, 215, 234, 253, 273, 293, 310, 330, 353, 374, 392, 413, 432, 453, 477, 499, 517, 538, 557, 578, 597, 619, 638, 653, 673, 693, 717, 742, 770, 799, 821, 844, 869, 892, 917, 942, 971, 1001, 1034, 1068, 1095, 1123, 1153, 1181, 1207, 1235, 1268, 1291, 1313, 1341, 1368, 1389, 1415, 1435, 1459, 1479, 1494, 1516, 1539}
)

func (i TagType) String() string {
	switch {
	case 61 <= i && i <= 64:
		i -= 61
		return _TagType_name_0[_TagType_index_0[i]:_TagType_index_0[i+1]]
	case i == 100:
		return _TagType_name_1
	case 256 <= i && i <= 271:
		i -= 256
		return _TagType_name_2[_TagType_index_2[i]:_TagType_index_2[i+1]]
	case i == 273:
		return _TagType_name_3
	case 1000 <= i && i <= 1069:
		i -= 1000
		return _TagType_name_4[_TagType_index_4[i]:_TagType_index_4[i+1]]
	case 1079 <= i && i <= 1196:
		i -= 1079
		return _TagType_name_5[_TagType_index_5[i]:_TagType_index_5[i+1]]
	case 5000 <= i && i <= 5027:
		i -= 5000
		return _TagType_name_6[_TagType_index_6[i]:_TagType_index_6[i+1]]
	case 5029 <= i && i <= 5097:
		i -= 5029
		return _TagType_name_7[_TagType_index_7[i]:_TagType_index_7[i+1]]
	default:
		return "TagType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
