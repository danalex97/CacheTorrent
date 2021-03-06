bia90 = [1444291.5, 1408389.2, 1387873.4, 1393146.2, 1388357.0]
bia50 = [1428534.75, 1393689.2, 1372269.6, 1376338.6, 1371253.4]
red_bia = [10.582703673873734, 12.52384302079788, 14.487532239664787, 16.899871746908353, 18.95719952865244]

bit90 = [1421906.0] * 5
bit50 = [1390034.5] * 5
red_bit = [46.89975] * 5

cac90 = [2795154.6, 2264567.0, 2075675.0, 1900072.0, 1764484.2]
cac50 = [2088354.6, 1712567.0, 1631675.0, 1545172.0, 1421284.2]
cacavg = [2174415.2, 1754970.4, 1633130.4, 1575467.7, 1506736.4]
red_cac = [4.369384835479256, 9.310080011430205, 14.231714285714286, 19.28725882460992, 24.299842857142856]

scac90  = [2465663.4, 1943782.8, 1718132.6, 1497803.8, 1383119.8]
scac50  = [1789485.6, 1382782.8, 1254632.6, 1183403.8, 1112819.8]
scacavg = [1822004.1, 1377073.8, 1228353.1, 1145023.0, 1127055.5]
red_scac = [4.135685973400614, 7.3762265684246815, 9.225472757071612, 10.801587820403821, 13.72113847664711]

sbia90  = [1283063.8, 1260930.5, 1252211.6, 1244652.5, 1244943.6]
sbia50  = [1264162.0, 1243663.7, 1235390.8, 1227402.5, 1227938.0]
sbiaavg = [1243139.3, 1224210.7, 1215837.3, 1209452.0, 1209937.8]
red_sbia = [11.08978937342432, 13.988364442537366, 14.439521789585445, 15.932816840853652, 17.56879524217716]

sbit90  = [1413451.2] * 5
sbit50  = [1406051.2] * 5
sbitavg = [1404587.0] * 5
red_sbit = [46.32947142857143] * 5

acac90 = [1311858.6, 1269271.8, 1261899.4, 1268813.4, 1294926.0]
acac50 = [1261158.6, 1264171.8, 1257930.6, 1265213.4, 1291426.0]
abia90 = [1272661.8, 1280000.0, 1283000.0, 1280051.6]
abia50 = [1268404.4, 1267638.0, 1280000.0, 1284261.5]
red_abia = [14.371836096967177, 16.073892809320004, 17.98625134264232, 20.011821596975324]
red_acac = [9.568944996001083, 13.641038166719099, 18.308232016698877, 24.228285714285715]

# Asym. big seed
sacac90 = [820490.6, 717618.0, 662806.5, 653595.0]
sacac50 = [609990.6, 554118.0, 536056.5, 512595.0]
sacacavg = [600125.8, 525791.9, 487413.7, 475073.5]
red_sacac = [7.296506708965068, 11.170763380130067, 15.82682142857143, 21.0987]

sabia90 = [626446.6, 632804.5, 643868.3, 634731.0]
sabia50 = [495946.6, 485804.5, 527868.3, 509106.0]
sabiaavg = [483029.2, 490061.9, 506963.9, 496343.3]
red_sabia = [11.688949840165073, 15.105469430672786, 17.38205260239677,  19.994200748891423]

sabit90  = [633349.0] * 4
sabit50  = [514099.0] * 4
sabitavg = [493678.3642857143] * 4
red_sabit = [47.1781428571428] * 4

idx = [10, 20, 30, 40, 50]

pic90 = [2674485.6, 2156115.6, 1803175.2, 1754077.2, 1779010.0]
pic50 = [2329785.6, 1880715.6, 1461679.8, 1428877.2, 1408510.0]

multi2_90  = [2290000.0, 1810500.0, 1658800.0, 1595500.0, 1565375.0]
multi2_50  = [1733125.0, 1469500.0, 1383700.0, 1348500.0, 1318250.0]
multi2_avg = [1812855.7, 1494886.4, 1399024.7, 1364757.1, 1346980.714285714]
multi2_red = [4.350125405908605, 7.352818223986509, 8.715263940917538, 9.834774046871459, 11.580178429992902]

multi3_90  = [2293250.0, 1785333.3, 1665750.0, 1677333.3, 1731833.3]
multi3_50  = [1739000.0, 1507833.3, 1436250.0, 1425833.3, 1478455.7]
multi3_avg = [1819113.2, 1522245.7, 1448213.7, 1449160.9, 1456833.3]
multi3_red = [4.037169839838482, 6.683063728897484, 8.11444742916304, 9.258826708650659, 10.216138696101217]

from plot import *

plots = [
	# # Intro
	# TimePlot(type="Bias")
	# 	.plot(idx, bia90, color='black', linestyle='-', label="Biased selection 90th percentile")
	# 	.plot(idx, bia50, color='black', linestyle='--', label="Biased selection 50th percentile")
	# 	.plot(idx, bit90, color='black', linestyle='-.', label="BitTorrent 90th percentile")
	# 	.plot(idx, bit50, color='black', linestyle=':', label="BitTorrent 50th percentile")
	# 	.save("plots/timeBias.png"),
	# TimePlot(type="Cache")
	# 	.plot(idx, cac90, color='black', linestyle='-', label="CacheTorrent 90th percentile")
	# 	.plot(idx, cac50, color='black', linestyle='--', label="CacheTorrent 50th percentile")
	# 	.plot(idx, bit90, color='black', linestyle='-.', label="BitTorrent 90th percentile")
	# 	.save("plots/timeSlowSeed.png"),
	# RedPlot(type="Cache")
	# 	.plot(idx, red_bit, color='black', linestyle='--', label="BitTorrent")
	# 	.plot(idx, red_cac, color='black', linestyle='-', label="CacheTorrent")
	# 	.save("plots/red.png"),
	# TimePlot(type="Cache")
	# 	.plot(idx, cacavg, color='black', linestyle='-', label="CacheTorrent average")
	# 	.plot(idx, cac50, color='black', linestyle='--', label="CacheTorrent 50th percentile")
	# 	.save("plots/cactime2.png"),
	#
	# # Big seed
	# TimePlot(type="Cache")
	# 	.plot(idx, scac90, color='black', linestyle='-', label="CacheTorrent 90th percentile - 8Mbps seed")
	# 	.plot(idx, scac50, color='black', linestyle=':', label="CacheTorrent 50th percentile - 8Mbps seed")
	# 	.plot(idx, cac90, color='black', linestyle='--', label="CacheTorrent 90th percentile - 2Mbps seed")
	# 	.plot(idx, cac50, color='black', linestyle='-.', label="CacheTorrent 50th percentile - 2Mbps seed")
	# 	.save("plots/scactime.png"),
	# RedPlot(type="Cache")
	# 	.plot(idx, red_scac, color='black', linestyle='-', label="CacheTorrent - 8Mbps seed")
	# 	.plot(idx, red_cac, color='black', linestyle='--', label="CacheTorrent - 2Mbps seed")
	# 	.save("plots/scacred.png"),
	#
	# # Biased
	# TimePlot(type="Both")
	# 	.plot(idx, scac50, color='black', linestyle='-', label="CacheTorrent 50th percentile")
	# 	.plot(idx, sbia50, color='black', linestyle='--', label="Biased selection 50th percentile")
	# 	.plot(idx, sbit50, color='black', linestyle=':', label="BitTorrent 50th percentile")
	# 	.save("plots/biaCache.png"),
	# TimePlot(type="Both")
	# 	.plot(idx, scac90, color='black', linestyle='-', label="CacheTorrent 90th percentile")
	# 	.plot(idx, sbia90, color='black', linestyle='--', label="Biased selection 90th percentile")
	# 	.plot(idx, sbit90, color='black', linestyle=':', label="BitTorrent 90th percentile")
	# 	.save("plots/biaCache90.png"),
	# RedPlot(type="Both")
	# 	.plot(idx, red_scac, color='black', linestyle='-', label="CacheTorrent")
	# 	.plot(idx, red_sbia, color='black', linestyle='--', label="Biased selection")
	# 	.plot(idx, red_sbit, color='black', linestyle=':', label="BitTorrent")
	# 	.setLoc(1)
	# 	.save("plots/cacred.png"),
	#
	# # Hetero
	# TimePlot(type="Both")
	# 	.plot(idx[1:], sabiaavg, color='black', linestyle=':', label="Biased selection average")
	# 	.plot(idx[1:], sacacavg, color='black', linestyle='-', label="CacheTorrent average")
	# 	.plot(idx[1:], sabitavg, color='black', linestyle='--', label="BitTorrent average")
	# 	.save("plots/asymAvg.png"),
	# RedPlot(type="Both")
	# 	.plot(idx[1:], red_sabia, color='black', linestyle=':', label="Biased selection")
	# 	.plot(idx[1:], red_sacac, color='black', linestyle='-', label="CacheTorrent")
	# 	.save("plots/asymRed.png"),
	# TimePlot(type="Cache")
	# 	.plot(idx[1:], sacac90, color='black', linestyle='--', label="CacheTorrent 90th percentile")
	# 	.plot(idx[1:], sacac50, color='black', linestyle='-', label="CacheTorrent 50th percentile")
	# 	.save("plots/asymPer.png"),
	#
	# TimePlot(type="Both")
	# 	.plot(idx[1:], abia50, color='black', linestyle='-', label="Biased selection 50th percentile")
	# 	.plot([20, 30, 35, 40, 50], acac90, color='black', linestyle='-.', label="CacheTorrent 90th percentile")
	# 	.plot([20, 30, 35, 40, 50], acac50, color='black', linestyle=':', label="CacheTorrent 50th percentile")
	# 	.save("plots/biaCacheAsym.png"),
	# RedPlot(type="Both")
	# 	.plot(idx, red_bia, color='black', linestyle='--', label="Biased selection")
	# 	.plot(idx, red_cac, color='black', linestyle='-', label="CacheTorrent")
	# 	.save("plots/biaCacheAsymRed.png"),
	#
	# # Pieces
	# PiecePlot()
	# 	.plot([i * 20 for i in idx], pic50, color='black', linestyle='-', label="50th percentile")
	# 	.plot([i * 20 for i in idx], pic90, color='black', linestyle='--', label="90th percentile")
	# 	.save("plots/pieces.png"),

	# MultiExtension
	TimePlot(type="Cache")
		.plot(idx[:4], scac90[:4], color='black', linestyle='-', label="CacheTorrent 90th percentile")
		.plot(idx[:4], multi2_90[:4], color='black', linestyle=':', label="MultiTorrent(2 sub-torrents) 90th percentile")
		.plot(idx[:4], multi3_90[:4], color='black', linestyle='--', label="MultiTorrent(3 sub-torrents) 90th percentile")
		.save("plots/multi90.png"),

	TimePlot(type="Cache")
		.plot(idx, multi2_50, color='black', linestyle='-.', label="MultiTorrent 50th percentile")
		.plot(idx, multi2_90, color='black', linestyle=':', label="MultiTorrent 90th percentile")
		.plot(idx, scac50, color='black', linestyle='-', label="CacheTorrent 50th percentile")
		.plot(idx, scac90, color='black', linestyle='--', label="CacheTorrent 90th percentile")
		.save("plots/multiDiff.png"),

	RedPlot(type="Cache")
		.plot(idx[:4], red_scac[:4], color='black', linestyle='-', label="CacheTorrent")
		.plot(idx[:4], multi2_red[:4], color='black', linestyle=':', label="MultiTorrent(2 sub-torrents)")
		.plot(idx[:4], multi3_red[:4], color='black', linestyle='--', label="MultiTorrent(3 sub-torrents)")
		.save("plots/multiRed.png"),
]
