HOST = $(shell hostname)

docker-ONU:
ifeq ($(HOST),dist106)
	docker build -t lab1:latest .
	docker rm -f onu
	docker run -it --name onu -p 50052:50052 --expose 50052 lab1:latest go run ONU/onu.go
else
	echo "Ejecutar SOLO en dist106"
endif

docker-continentes:
	docker build -t lab1:latest .
	docker rm -f regional
ifeq ($(HOST),localhost)
	docker run  -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Asia/regional.go
endif
ifeq ($(HOST),dist106)
	docker run  -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Europa/regional.go
endif
ifeq ($(HOST),dist107)
	docker run  -it --rm --name regional --expose 50052 lab1:latest go run Regionales/LatinoAmerica/regional.go
endif
ifeq ($(HOST),dist108)
	docker run  -it --rm --name regional --expose 50052 lab1:latest go run Regionales/Australia/regional.go
endif

docker-OMS:
ifeq ($(HOST),localhost)
	docker build -t lab1:latest .
	docker rm -f oms
	docker run  -it --name oms -p 50052:50052 -p 50053:50053 --expose 50052 --expose 50053 lab1:latest go run OMS/oms.go
else
	echo "Ejecutar SOLO en dist105"
endif

dockeratanode:
	docker build -t lab1:latest .
	docker rm -f datanode
ifeq ($(HOST),dist107)
	docker run  -it --name datanode -p 50052:50052 --expose 50052 lab1:latest go run DataNode/Data1/datanode.go
else ifeq ($(HOST),dist108)
	docker run  -it --name datanode -p 50052:50052 --expose 50052 lab1:latest go run DataNode/Data2/datanode.go
else
	echo "Ejecutar SOLO en dist107 y dist108"
endif

clean:
	rm OMS/DATA.txt
	touch OMS/DATA.txt

	rm DataNode/Data1/DATA.txt
	touch DataNode/Data1/DATA.txt

	rm DataNode/Data2/DATA.txt
	touch DataNode/Data2/DATA.txt
	