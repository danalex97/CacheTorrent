bia90 = [1444291.5, 1408389.2, 1387873.4, 1393146.2, 1388357.0]
bia50 = [1428534.75, 1393689.2, 1372269.6, 1376338.6, 1371253.4]
red_bia = [10.582703673873734, 12.52384302079788, 14.487532239664787, 16.899871746908353, 18.95719952865244]

bit90 = [1421906.0] * 5
bit50 = [1390034.5] * 5
red_bit = [46.89975] * 5

cac90 = [2795154.6, 2264567.0, 2075675.0, 1900072.0, 1764484.2]
cac50 = [2088354.6, 1712567.0, 1631675.0, 1545172.0, 1421284.2]
red_cac = [4.369384835479256, 9.310080011430205, 14.231714285714286, 19.28725882460992, 24.299842857142856]

scac90  = [2465663.4, 1943782.8, 1718132.6, 1497803.8, 1383119.8]
scac50  = [1789485.6, 1382782.8, 1254632.6, 1183403.8, 1112819.8]
scacavg = [2174415.2, 1754970.4, 1633130.4, 1575467.7, 1506736.4]
red_scac = [4.135685973400614, 7.3762265684246815, 9.225472757071612, 10.801587820403821, 13.72113847664711]

sbit90  = [1413451.2] * 5
sbit50  = [1406051.2] * 5
sbitavg = [1404587.0] * 5

acac90 = [1311858.6, 1269271.8, 1261899.4, 1268813.4, 1294926.0]
acac50 = [1261158.6, 1264171.8, 1257930.6, 1265213.4, 1291426.0]
abia90 = [1272661.8, 1280000.0, 1283000.0, 1280051.6]
abia50 = [1268404.4, 1267638.0, 1280000.0, 1284261.5]
red_abia = [14.371836096967177, 16.073892809320004, 17.98625134264232, 20.011821596975324]
red_acac = [9.568944996001083, 13.641038166719099, 18.308232016698877, 24.228285714285715]

idx = [10, 20, 30, 40, 50]

pic90 = [2674485.6, 2156115.6, 1803175.2, 1754077.2, 1779010.0]
pic50 = [2329785.6, 1880715.6, 1461679.8, 1428877.2, 1408510.0]

from plot import *

plots = [
	TimePlot(type="Bias")
		.plot(idx, bia90, color='black', linestyle='-', label="Biased selection 90th percentile")
		.plot(idx, bia50, color='black', linestyle='--', label="Biased selection 50th percentile")
		.plot(idx, bit90, color='black', linestyle='-.', label="BitTorrent 90th percentile")
		.plot(idx, bit50, color='black', linestyle=':', label="BitTorrent 50th percentile")
		.save("plots/timeBias.png"),
	TimePlot(type="Cache")
		.plot(idx, cac90, color='black', linestyle='-', label="CacheTorrent 90th percentile")
		.plot(idx, cac50, color='black', linestyle='--', label="CacheTorrent 50th percentile")
		.plot(idx, bit90, color='black', linestyle='-.', label="BitTorrent 90th percentile")
		.save("plots/timeSlowSeed.png"),
	RedPlot(type="Cache")
		.plot(idx, red_bit, color='black', linestyle='--', label="BitTorrent")
		.plot(idx, red_cac, color='black', linestyle='-', label="CacheTorrent")
		.save("plots/red.png"),
	TimePlot(type="Cache")
		.plot(idx, scacavg, color='black', linestyle='-', label="CacheTorrent average")
		.plot(idx, cac50, color='black', linestyle='--', label="CacheTorrent 50th percentile")
		.save("plots/cactime2.png"),
	TimePlot(type="Cache")
		.plot(idx, scac90, color='black', linestyle='-', label="CacheTorrent 90th percentile - 8Mbps seed")
		.plot(idx, scac50, color='black', linestyle=':', label="CacheTorrent 50th percentile - 8Mbps seed")
		.plot(idx, cac90, color='black', linestyle='--', label="CacheTorrent 90th percentile - 2Mbps seed")
		.plot(idx, cac50, color='black', linestyle='-.', label="CacheTorrent 50th percentile - 2Mbps seed")
		.save("plots/scactime.png"),
	RedPlot(type="Cache")
		.plot(idx, red_scac, color='black', linestyle='-', label="CacheTorrent - 8Mbps seed")
		.plot(idx, red_cac, color='black', linestyle='--', label="CacheTorrent - 2Mbps seed")
		.save("plots/scacred.png"),
	TimePlot(type="Both")
		.plot(idx[1:], abia50, color='black', linestyle='-', label="Biased selection 50th percentile")
		.plot([20, 30, 35, 40, 50], acac90, color='black', linestyle='-.', label="CacheTorrent 90th percentile")
		.plot([20, 30, 35, 40, 50], acac50, color='black', linestyle=':', label="CacheTorrent 50th percentile")
		.save("plots/biaCacheAsym.png"),
	RedPlot(type="Both")
		.plot(idx, red_bia, color='black', linestyle='--', label="Biased selection")
		.plot(idx, red_cac, color='black', linestyle='-', label="CacheTorrent")
		.save("plots/biaCacheAsymRed.png"),
	PiecePlot()
		.plot([i * 20 for i in idx], pic50, color='black', linestyle='-', label="50th percentile")
		.plot([i * 20 for i in idx], pic90, color='black', linestyle='--', label="90th percentile")
		.save("plots/pieces.png"),
]
