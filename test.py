
import time

start = time.time()

list = []

i = 0.0

while True:
	i=i+1.0
	if i > 999999.0:
		break
	

end = time.time()
print((end-start)/1000)