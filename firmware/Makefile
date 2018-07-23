DEVICE?=GreenHouse
#DEVICE?=outside
DEVICE_TYPE?=photon

show:
	@echo flashing ${DEVICE}

firmware:
	@particle compile ${DEVICE_TYPE} photonWeather.ino homestead_store.cpp  homestead_store.h SparkFun_Photon_Weather_Shield_Library.cpp  SparkFun_Photon_Weather_Shield_Library.h

flash:
	@particle flash ${DEVICE} photonWeather.ino homestead_store.cpp  homestead_store.h SparkFun_Photon_Weather_Shield_Library.cpp  SparkFun_Photon_Weather_Shield_Library.h
