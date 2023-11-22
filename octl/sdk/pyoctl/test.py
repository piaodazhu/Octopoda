import pyoctl
import time
import os
import platform

libpath = "./libcoctl.so"
if platform.system() == "Windows":
    libpath = "./libcoctl.dll"

def main():
	octl = pyoctl.OctlClient(libpath, "../../octl_test.yaml")
	binfo, ninfos = octl.get_nodes_info_list(['pi4', 'pi5'])
	print(binfo)
	for info in ninfos:
		print(info)
		ninfo = octl.get_node_info(info.name)
		if info.name != ninfo.name:
			raise Exception("get node info failed.")
	
	status_list = octl.get_nodes_status_list([])
	for status in status_list:
		print(status)
		nstatus = octl.get_node_status(status.name)
		if status.name != nstatus.name:
			raise Exception("get node status failed.")
	
	results = octl.run(r"{ifconfig}", ['pi4', 'pi5'], True)
	for result in results:
		print(result)
	
	results = octl.xrun(r"{ls}", ['pi4', 'pi5'], 1)
	for result in results:
		print(result)
	
	try:
		octl.del_group("setByPy")
	except pyoctl.OctlException as e:
		print(e)

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

	octl.prune_nodes()
	scens = octl.get_scenarios_info_list()
	for scen in scens:
		name = scen['Name']
		info = octl.get_scenario_info(name)
		version = octl.get_scenario_version(name)
		print(f"found scenario {name} with {info['Version']} (total version number is {len(version)})")

	nodeapps = octl.get_nodeapps_info_list("pi4")
	for nodeapp in nodeapps:
		name = nodeapp['Name']
		tmp = str(name).split("@")
		info = octl.get_nodeapp_info("pi4", tmp[0], tmp[1])
		print(f"found pi4 nodeapp {info['Name']} with {info['Versions']}")
	
	octl.apply_scenario("../../s1", "start", "byPyoctl", 0)
	print("PASS ALL")


if __name__ == "__main__":
	main()