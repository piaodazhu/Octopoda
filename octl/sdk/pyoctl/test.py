import pyoctl
import time
import os

def main():
	octl = pyoctl.OctlClient("./libcoctl.so", "../../octl_test.yaml")
	binfo, ninfos = octl.get_nodes_info_list(['pi4', 'pi5'])
	print(binfo)
	for info in ninfos:
		print(info)
	
	status_list = octl.get_nodes_status_list([])
	for status in status_list:
		print(status)
	
	# results = octl.run(r"{ifconfig}", ['pi4', 'pi5'])
	# for result in results:
	# 	print(result)
	
	results = octl.xrun(r"{ls}", ['pi4', 'pi5'], 1)
	for result in results:
		print(result)
	
	try:
		octl.del_group("setByPy")
	except Exception:
		pass
	
	octl.set_group("setByPy", False, ['pi4', 'pi5', 'yang'])

	group_names = octl.get_groups_list()
	print(group_names)
	for group in group_names:
		members = octl.get_group(group)
		print("members of group " + group + ": ", members)
	
	octl.del_group("setByPy")

	time.sleep(1)
	print(octl.get_groups_list())

	file_content = "hello world"
	file = open("./testfile", "w+")
	file.write(file_content)
	file.close()

	results = octl.distribute_file("./testfile", "distributeByPy/", ['pi4', 'pi5'])
	for result in results:
		print(result)
	
	os.remove("./testfile")

	result = octl.pull_file("fstore", "pi5", "distributeByPy/testfile", "fromPi5")
	print(result)

	file = open("fromPi5/testfile", "r")
	if file.read() != file_content:
		raise Exception("file distribute or pull failed.")
	file.close()

	os.remove("fromPi5/testfile")
	os.rmdir("fromPi5")

	print("PASS ALL")


if __name__ == "__main__":
	main()