HOST = $(shell hostname)

docker-ONU:
ifeq ($(HOST),dist106)
	docker build -t lab1:latest .
	docker rm -f onu
	docker run -d -it --name onu -p 50052:50052 --expose 50052 lab1:latest go run ONU/onu.go
else
	echo "Ejecutar SOLO en dist106"
endif

docker-continentes:
	docker build -t lab1:latest .
	docker rm -f regional
ifeq ($(HOST),localhost)
	docker run -d -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Asia/regional.go
endif
ifeq ($(HOST),dist106)
	docker run -d -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Europa/regional.go
endif
ifeq ($(HOST),dist107)
	docker run -d -it --rm --name regional --expose 50052 lab1:latest go run Regionales/LatinoAmerica/regional.go
endif
ifeq ($(HOST),dist108)
	docker run -d -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Australia/regional.go
endif

docker-OMS:
ifeq ($(HOST),localhost)
	docker build -t lab1:latest .
	docker rm -f oms
	rm -rf OMS/DATA.txt
	docker run -d -it --name oms -p 50052:50052 -p 50053:50053 --expose 50052 --expose 50053 lab1:latest go run OMS/oms.go
else
	echo "Ejecutar SOLO en dist105"
endif

docker-datanode:
	docker build -t lab1:latest .
	docker rm -f datanode
	rm -rf OMS/DATA.txt
ifeq ($(HOST),dist107)
	docker run -d -it --name datanode -p 50052:50052 --expose 50052 lab1:latest go run DataNode/Data1/datanode.go
endif
ifeq ($(HOST),dist108)
	docker run -d -it --name datanode -p 50052:50052 --expose 50052 lab1:latest go run DataNode/Data2/datanode.go
else
	echo "Ejecutar SOLO en dist107 y dist108"
endif