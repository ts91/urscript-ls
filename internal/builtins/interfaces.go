package builtins

func interfacesBuiltins() []Builtin {
	return []Builtin{
		// External F/T sensor
		{
			Name:      "enable_external_ft_sensor",
			Signature: "enable_external_ft_sensor(enable, sensor_mass=0.0, sensor_measuring_offset=[0.0,0.0,0.0], sensor_cog=[0.0,0.0,0.0])",
			Module:    "interfaces",
			Doc: "Deprecated: Enable external F/T sensor input via RTDE. Use ft_rtde_input_enable instead.\n" +
				"@param enable bool Enable or disable\n" +
				"@param sensor_mass float Sensor mass [kg]\n" +
				"@param sensor_measuring_offset list Measuring offset [x,y,z] in meters from tool flange\n" +
				"@param sensor_cog list Center of gravity [x,y,z] in meters from tool flange",
		},
		{
			Name:      "ft_rtde_input_enable",
			Signature: "ft_rtde_input_enable(enable, sensor_mass=0.0, sensor_measuring_offset=[0.0,0.0,0.0], sensor_cog=[0.0,0.0,0.0])",
			Module:    "interfaces",
			Doc: "Enable or disable external F/T sensor input via RTDE.\n" +
				"Affects force_mode, screw_driving, and freedrive_mode.\n" +
				"@param enable bool Enable or disable\n" +
				"@param sensor_mass float Sensor mass [kg]\n" +
				"@param sensor_measuring_offset list Measuring offset [x,y,z] in meters from tool flange\n" +
				"@param sensor_cog list Center of gravity [x,y,z] in meters from tool flange",
		},
		// Analog I/O
		{
			Name:      "get_analog_in",
			Signature: "get_analog_in(n)",
			Module:    "interfaces",
			Doc: "Deprecated: Read analog input signal level. Use get_standard_analog_in or get_tool_analog_in.\n" +
				"@param n int Input number [0-3] (ports 2-3 map to tool inputs)\n" +
				"@returns float Signal level in Amperes or Volts",
		},
		{
			Name:      "get_analog_out",
			Signature: "get_analog_out(n)",
			Module:    "interfaces",
			Doc: "Deprecated: Read analog output signal level. Use get_standard_analog_out.\n" +
				"@param n int Output number [0-1]\n" +
				"@returns float Signal level in Amperes or Volts",
		},
		{
			Name:      "get_standard_analog_in",
			Signature: "get_standard_analog_in(n)",
			Module:    "interfaces",
			Doc: "Read standard analog input signal level.\n" +
				"@param n int Input number [0-1]\n" +
				"@returns float Signal level in Amperes or Volts",
		},
		{
			Name:      "get_standard_analog_out",
			Signature: "get_standard_analog_out(n)",
			Module:    "interfaces",
			Doc: "Read standard analog output signal level.\n" +
				"@param n int Output number [0-1]\n" +
				"@returns float Signal level in Amperes or Volts",
		},
		{
			Name:      "get_tool_analog_in",
			Signature: "get_tool_analog_in(n)",
			Module:    "interfaces",
			Doc: "Read tool analog input signal level.\n" +
				"@param n int Input number [0-1]\n" +
				"@returns float Signal level in Amperes or Volts",
		},
		{
			Name:      "set_analog_out",
			Signature: "set_analog_out(n, f)",
			Module:    "interfaces",
			Doc: "Deprecated: Set analog output signal level. Use set_standard_analog_out.\n" +
				"@param n int Output number [0-1]\n" +
				"@param f float Relative signal level [0-1]",
		},
		{
			Name:      "set_standard_analog_out",
			Signature: "set_standard_analog_out(n, f)",
			Module:    "interfaces",
			Doc: "Set standard analog output signal level.\n" +
				"@param n int Output number [0-1]\n" +
				"@param f float Relative signal level [0-1]",
		},
		// Digital I/O
		{
			Name:      "get_configurable_digital_in",
			Signature: "get_configurable_digital_in(n)",
			Module:    "interfaces",
			Doc: "Read configurable digital input signal level.\n" +
				"@param n int Input number [0-7]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_configurable_digital_out",
			Signature: "get_configurable_digital_out(n)",
			Module:    "interfaces",
			Doc: "Read configurable digital output signal level.\n" +
				"@param n int Output number [0-7]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_digital_in",
			Signature: "get_digital_in(n)",
			Module:    "interfaces",
			Doc: "Deprecated: Read digital input. Use get_standard_digital_in or get_tool_digital_in.\n" +
				"@param n int Input number [0-9] (ports 8-9 map to tool inputs)\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_digital_out",
			Signature: "get_digital_out(n)",
			Module:    "interfaces",
			Doc: "Deprecated: Read digital output. Use get_standard_digital_out or get_tool_digital_out.\n" +
				"@param n int Output number [0-9] (ports 8-9 map to tool outputs)\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_standard_digital_in",
			Signature: "get_standard_digital_in(n)",
			Module:    "interfaces",
			Doc: "Read standard digital input signal level.\n" +
				"@param n int Input number [0-7]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_standard_digital_out",
			Signature: "get_standard_digital_out(n)",
			Module:    "interfaces",
			Doc: "Read standard digital output signal level.\n" +
				"@param n int Output number [0-7]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_tool_digital_in",
			Signature: "get_tool_digital_in(n)",
			Module:    "interfaces",
			Doc: "Read tool digital input signal level.\n" +
				"@param n int Input number [0-1]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "get_tool_digital_out",
			Signature: "get_tool_digital_out(n)",
			Module:    "interfaces",
			Doc: "Read tool digital output signal level.\n" +
				"@param n int Output number [0-1]\n" +
				"@returns bool Signal level",
		},
		{
			Name:      "set_configurable_digital_out",
			Signature: "set_configurable_digital_out(n, b)",
			Module:    "interfaces",
			Doc: "Set configurable digital output signal level.\n" +
				"@param n int Output number [0-7]\n" +
				"@param b bool Signal level",
		},
		{
			Name:      "set_digital_out",
			Signature: "set_digital_out(n, b)",
			Module:    "interfaces",
			Doc: "Deprecated: Set digital output. Use set_standard_digital_out or set_tool_digital_out.\n" +
				"@param n int Output number [0-9] (ports 8-9 map to tool outputs)\n" +
				"@param b bool Signal level",
		},
		{
			Name:      "set_standard_digital_out",
			Signature: "set_standard_digital_out(n, b)",
			Module:    "interfaces",
			Doc: "Set standard digital output signal level.\n" +
				"@param n int Output number [0-7]\n" +
				"@param b bool Signal level",
		},
		{
			Name:      "set_tool_digital_out",
			Signature: "set_tool_digital_out(n, b)",
			Module:    "interfaces",
			Doc: "Set tool digital output signal level.\n" +
				"@param n int Output number [0-1]\n" +
				"@param b bool Signal level",
		},
		// Flags
		{
			Name:      "get_flag",
			Signature: "get_flag(n)",
			Module:    "interfaces",
			Doc: "Read an internal flag (persists between program runs).\n" +
				"@param n int Flag number [0-31]\n" +
				"@returns bool Stored value",
		},
		{
			Name:      "set_flag",
			Signature: "set_flag(n, b)",
			Module:    "interfaces",
			Doc: "Set an internal flag (persists between program runs).\n" +
				"@param n int Flag number [0-31]\n" +
				"@param b bool Value to store",
		},
		// Modbus
		{
			Name:      "modbus_add_signal",
			Signature: "modbus_add_signal(IP, slave_number, signal_address, signal_type, signal_name, sequential_mode=False, register_count=1)",
			Module:    "interfaces",
			Doc: "Add a Modbus signal for the controller to supervise.\n" +
				"@param IP string IP address of the Modbus unit\n" +
				"@param slave_number int Slave number [0-255]\n" +
				"@param signal_address int Address of the coil or register\n" +
				"@param signal_type int Signal type (0=digital in, 1=digital out, 2=register in, 3=register out, 15=multi digital out, 16=multi register out)\n" +
				"@param signal_name string Unique name for the signal (max 20 chars)\n" +
				"@param sequential_mode bool Force waiting for response before next request\n" +
				"@param register_count int Number of registers/coils [1-123]",
		},
		{
			Name:      "modbus_add_rw_signal",
			Signature: "modbus_add_rw_signal(IP, slave_number, read_address, read_register_count, write_address, write_register_count, signal_name, sequential_mode=False)",
			Module:    "interfaces",
			Doc: "Add a Modbus read-write signal using function code 23.\n" +
				"@param IP string IP address of the Modbus unit\n" +
				"@param slave_number int Slave number [0-255]\n" +
				"@param read_address int First register address to read from\n" +
				"@param read_register_count int Number of registers to read [1-123]\n" +
				"@param write_address int First register address to write to\n" +
				"@param write_register_count int Number of registers to write [1-123]\n" +
				"@param signal_name string Unique signal name (max 20 chars)\n" +
				"@param sequential_mode bool Force waiting for response before next request",
		},
		{
			Name:      "modbus_delete_signal",
			Signature: "modbus_delete_signal(signal_name)",
			Module:    "interfaces",
			Doc: "Delete a previously added Modbus signal.\n" +
				"@param signal_name string Name of the signal to delete",
		},
		{
			Name:      "modbus_get_signal_status",
			Signature: "modbus_get_signal_status(signal_name, is_secondary_program=False)",
			Module:    "interfaces",
			Doc: "Read the current value of a Modbus signal.\n" +
				"@param signal_name string Name of the signal\n" +
				"@param is_secondary_program bool Internal use, must be False\n" +
				"@returns any Bool for digital signals, int for registers (or list for multi-register signals)",
		},
		{
			Name:      "modbus_send_custom_command",
			Signature: "modbus_send_custom_command(IP, slave_number, function_code, data)",
			Module:    "interfaces",
			Doc: "Send a custom Modbus command. Cannot request data (response is not received).\n" +
				"@param IP string IP address of the Modbus unit\n" +
				"@param slave_number int Slave number [0-255]\n" +
				"@param function_code int Modbus function code\n" +
				"@param data list Array of byte values [0-255]",
		},
		{
			Name:      "modbus_set_digital_input_action",
			Signature: "modbus_set_digital_input_action(signal_name, action)",
			Module:    "interfaces",
			Doc: "Set a Modbus digital input signal to a default or freedrive action.\n" +
				"@param signal_name string Name of a previously added digital input signal\n" +
				"@param action string Action type: \"default\" or \"freedrive\"",
		},
		{
			Name:      "modbus_set_output_register",
			Signature: "modbus_set_output_register(signal_name, register_value, is_secondary_program=False)",
			Module:    "interfaces",
			Doc: "Set the value of a Modbus output register signal.\n" +
				"@param signal_name string Name of an output register signal\n" +
				"@param register_value any Integer [0-65535] or list of integers\n" +
				"@param is_secondary_program bool Internal use, must be False",
		},
		{
			Name:      "modbus_set_output_signal",
			Signature: "modbus_set_output_signal(signal_name, digital_value, is_secondary_program=False)",
			Module:    "interfaces",
			Doc: "Set the value of a Modbus digital output signal.\n" +
				"@param signal_name string Name of an output digital signal\n" +
				"@param digital_value any Bool or list of bools\n" +
				"@param is_secondary_program bool Internal use, must be False",
		},
		{
			Name:      "modbus_set_signal_update_frequency",
			Signature: "modbus_set_signal_update_frequency(signal_name, update_frequency)",
			Module:    "interfaces",
			Doc: "Set the polling frequency for a Modbus signal.\n" +
				"@param signal_name string Name of the signal\n" +
				"@param update_frequency int Update frequency [0-500] Hz (-1 or 0 for acyclic)",
		},
		{
			Name:      "modbus_get_error",
			Signature: "modbus_get_error(signal_name)",
			Module:    "interfaces",
			Doc: "Returns the current error state of a Modbus signal.\n" +
				"@param signal_name string Name of the signal\n" +
				"@returns int Error code (0=no error, see Modbus device/connection error codes)",
		},
		{
			Name:      "modbus_get_time_since_signal_invalid",
			Signature: "modbus_get_time_since_signal_invalid(signal_name)",
			Module:    "interfaces",
			Doc: "Returns the time since the signal last had a communication or device error.\n" +
				"@param signal_name string Name of the signal\n" +
				"@returns float Seconds since the signal became invalid",
		},
		{
			Name:      "modbus_request_update_signal_value",
			Signature: "modbus_request_update_signal_value(signal_name)",
			Module:    "interfaces",
			Doc: "Request an immediate update of a Modbus signal regardless of polling frequency.\n" +
				"@param signal_name string Name of the signal",
		},
		{
			Name:      "modbus_reset_connection",
			Signature: "modbus_reset_connection(connection_id, is_blocking=True)",
			Module:    "interfaces",
			Doc: "Tear down and reconnect all signals to a Modbus device.\n" +
				"@param connection_id string IP address of the Modbus unit\n" +
				"@param is_blocking bool Block until connection is restored (default True)",
		},
		{
			Name:      "modbus_set_signal_watchdog",
			Signature: "modbus_set_signal_watchdog(signal_name, new_timeout_in_sec)",
			Module:    "interfaces",
			Doc: "Set error tolerance timeout for a Modbus signal (default is 2 seconds).\n" +
				"@param signal_name string Name of the signal\n" +
				"@param new_timeout_in_sec float Timeout in seconds [0-300]",
		},
		// Registers
		{
			Name:      "read_input_boolean_register",
			Signature: "read_input_boolean_register(address)",
			Module:    "interfaces",
			Doc: "Read a boolean from an input register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-127]\n" +
				"@returns bool Register value",
		},
		{
			Name:      "read_input_float_register",
			Signature: "read_input_float_register(address)",
			Module:    "interfaces",
			Doc: "Read a float from an input register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@returns float Register value",
		},
		{
			Name:      "read_input_integer_register",
			Signature: "read_input_integer_register(address)",
			Module:    "interfaces",
			Doc: "Read an integer from an input register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@returns int Register value",
		},
		{
			Name:      "read_output_boolean_register",
			Signature: "read_output_boolean_register(address)",
			Module:    "interfaces",
			Doc: "Read a boolean from an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-127]\n" +
				"@returns bool Register value",
		},
		{
			Name:      "read_output_float_register",
			Signature: "read_output_float_register(address)",
			Module:    "interfaces",
			Doc: "Read a float from an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@returns float Register value",
		},
		{
			Name:      "read_output_integer_register",
			Signature: "read_output_integer_register(address)",
			Module:    "interfaces",
			Doc: "Read an integer from an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@returns int Register value",
		},
		{
			Name:      "read_port_bit",
			Signature: "read_port_bit(address)",
			Module:    "interfaces",
			Doc: "Read a port bit accessible by Modbus clients.\n" +
				"@param address int Port address\n" +
				"@returns bool Port value",
		},
		{
			Name:      "read_port_register",
			Signature: "read_port_register(address)",
			Module:    "interfaces",
			Doc: "Read a port register accessible by Modbus clients.\n" +
				"@param address int Port address\n" +
				"@returns int Signed integer value [-32768, 32767]",
		},
		{
			Name:      "write_output_boolean_register",
			Signature: "write_output_boolean_register(address, value)",
			Module:    "interfaces",
			Doc: "Write a boolean to an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-127]\n" +
				"@param value bool Value to write",
		},
		{
			Name:      "write_output_float_register",
			Signature: "write_output_float_register(address, value)",
			Module:    "interfaces",
			Doc: "Write a float to an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@param value float Value to write",
		},
		{
			Name:      "write_output_integer_register",
			Signature: "write_output_integer_register(address, value)",
			Module:    "interfaces",
			Doc: "Write an integer to an output register (accessible via FieldBus/RTDE).\n" +
				"@param address int Register address [0-47]\n" +
				"@param value int Value to write",
		},
		{
			Name:      "write_port_bit",
			Signature: "write_port_bit(address, value)",
			Module:    "interfaces",
			Doc: "Write a port bit accessible by Modbus clients.\n" +
				"@param address int Port address\n" +
				"@param value bool Value to set",
		},
		{
			Name:      "write_port_register",
			Signature: "write_port_register(address, value)",
			Module:    "interfaces",
			Doc: "Write a port register accessible by Modbus clients.\n" +
				"@param address int Port address\n" +
				"@param value int Value to set [0-65535] or [-32768, 32767]",
		},
		// RPC
		{
			Name:      "rpc_factory",
			Signature: "rpc_factory(type, url)",
			Module:    "interfaces",
			Doc: "Create a Remote Procedure Call handle for calling functions on a remote server.\n" +
				"@param type string RPC backend type (currently only \"xmlrpc\")\n" +
				"@param url string Server URL (pstream: \"ip:port\", http: \"http://ip:port/path\")\n" +
				"@returns handle RPC handle for calling remote functions",
		},
		// RTDE
		{
			Name:      "rtde_set_watchdog",
			Signature: "rtde_set_watchdog(variable_name, min_frequency, action=\"pause\")",
			Module:    "interfaces",
			Doc: "Set a watchdog on an RTDE input variable. Triggers action if updates stop.\n" +
				"@param variable_name string RTDE input variable name\n" +
				"@param min_frequency float Minimum expected update frequency [Hz]\n" +
				"@param action string Action on violation: \"ignore\", \"pause\", or \"stop\"",
		},
		// Analog input range
		{
			Name:      "set_analog_inputrange",
			Signature: "set_analog_inputrange(port, range)",
			Module:    "interfaces",
			Doc: "Deprecated: Set analog input range. Use set_standard_analog_input_domain or set_tool_analog_input_domain.\n" +
				"@param port int Analog input port (0-1=controller, 2-3=tool)\n" +
				"@param range int Range setting (controller: 0=0-5V, 2=0-10V; tool: 0=0-5V, 1=0-10V, 2=4-20mA)",
		},
		// Tool communication
		{
			Name:      "set_tool_communication",
			Signature: "set_tool_communication(enabled, baud_rate, parity, stop_bits, rx_idle_chars, tx_idle_chars)",
			Module:    "interfaces",
			Doc: "Enable or disable Tool Communication Interface (TCI). Disables tool analog inputs when enabled.\n" +
				"@param enabled bool Enable or disable TCI\n" +
				"@param baud_rate int Baud rate (9600, 19200, 38400, 57600, 115200, 1000000, 2000000, 5000000)\n" +
				"@param parity int Parity (0=none, 1=odd, 2=even)\n" +
				"@param stop_bits int Stop bits (1 or 2)\n" +
				"@param rx_idle_chars float RX idle chars before message complete [1.0-40.0]\n" +
				"@param tx_idle_chars float TX idle chars before new transmission [0.0-40.0]",
		},
		{
			Name:      "set_tool_digital_output_mode",
			Signature: "set_tool_digital_output_mode(n, mode)",
			Module:    "interfaces",
			Doc: "Set the output mode of a tool digital output pin.\n" +
				"@param n int Output pin number [0-1]\n" +
				"@param mode int Pin mode (1=Sinking/NPN, 2=Sourcing/PNP, 3=Push-Pull)",
		},
		{
			Name:      "set_tool_output_mode",
			Signature: "set_tool_output_mode(mode)",
			Module:    "interfaces",
			Doc: "Set the tool digital output mode.\n" +
				"@param mode int 1=power (dual pin) mode, outputs used as extra supply",
		},
		{
			Name:      "set_tool_voltage",
			Signature: "set_tool_voltage(voltage)",
			Module:    "interfaces",
			Doc: "Set the tool connector power supply voltage.\n" +
				"@param voltage int Voltage level: 0, 12, or 24 [V]",
		},
		// Sockets
		{
			Name:      "socket_close",
			Signature: "socket_close(socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Close a TCP/IP socket connection.\n" +
				"@param socket_name string Name of the socket",
		},
		{
			Name:      "socket_get_var",
			Signature: "socket_get_var(name, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send \"GET <name>\\n\" and read an integer response from the server. Returns 0 on timeout (2s).\n" +
				"@param name string Variable name\n" +
				"@param socket_name string Name of the socket\n" +
				"@returns int Value from the server, 0 on timeout",
		},
		{
			Name:      "socket_open",
			Signature: "socket_open(address, port, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Open a TCP/IP socket connection. Times out after 2 seconds.\n" +
				"@param address string Server address\n" +
				"@param port int Port number\n" +
				"@param socket_name string Name for the socket\n" +
				"@returns bool True if connection established, False on failure",
		},
		{
			Name:      "socket_read_ascii_float",
			Signature: "socket_read_ascii_float(number, socket_name=\"socket_0\", timeout=2)",
			Module:    "interfaces",
			Doc: "Read ASCII-formatted floats from the socket (max 30 values).\n" +
				"Returns [count, val1, val2, ...]. Returns [0, nan, ...] on timeout.\n" +
				"@param number int Number of floats to read\n" +
				"@param socket_name string Socket name\n" +
				"@param timeout float Timeout in seconds (0 or negative = wait indefinitely)\n" +
				"@returns list List of floats (length = number+1, first element is count)",
		},
		{
			Name:      "socket_read_binary_integer",
			Signature: "socket_read_binary_integer(number, socket_name=\"socket_0\", timeout=2)",
			Module:    "interfaces",
			Doc: "Read 32-bit integers in network byte order from the socket (max 30 values).\n" +
				"Returns [count, val1, val2, ...]. Returns [0, -1, ...] on timeout.\n" +
				"@param number int Number of integers to read\n" +
				"@param socket_name string Socket name\n" +
				"@param timeout float Timeout in seconds (0 or negative = wait indefinitely)\n" +
				"@returns list List of ints (length = number+1, first element is count)",
		},
		{
			Name:      "socket_read_byte_list",
			Signature: "socket_read_byte_list(number, socket_name=\"socket_0\", timeout=2)",
			Module:    "interfaces",
			Doc: "Read bytes from the socket (max 30 values).\n" +
				"Returns [count, b1, b2, ...]. Returns [0, -1, ...] on timeout.\n" +
				"@param number int Number of bytes to read\n" +
				"@param socket_name string Socket name\n" +
				"@param timeout float Timeout in seconds (0 or negative = wait indefinitely)\n" +
				"@returns list List of ints (length = number+1, first element is count)",
		},
		{
			Name:      "socket_read_line",
			Signature: "socket_read_line(socket_name=\"socket_0\", timeout=2)",
			Module:    "interfaces",
			Doc: "Deprecated: Read until newline from socket. Use socket_read_string with interpret_escape=True.\n" +
				"@param socket_name string Socket name\n" +
				"@param timeout float Timeout in seconds\n" +
				"@returns string Line without newline characters, or empty string on timeout",
		},
		{
			Name:      "socket_read_string",
			Signature: "socket_read_string(socket_name=\"socket_0\", prefix=\"\", suffix=\"\", interpret_escape=False, timeout=2)",
			Module:    "interfaces",
			Doc: "Read a string from the socket (max 1024 chars). Use prefix/suffix to extract substrings.\n" +
				"@param socket_name string Socket name\n" +
				"@param prefix string Start delimiter (data before and including prefix is discarded)\n" +
				"@param suffix string End delimiter (data after suffix is preserved on socket)\n" +
				"@param interpret_escape bool Allow \\n, \\r, \\t in prefix/suffix\n" +
				"@param timeout float Timeout in seconds (0 or negative = wait indefinitely)\n" +
				"@returns string Received string, or empty on timeout",
		},
		{
			Name:      "socket_send_byte",
			Signature: "socket_send_byte(value, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send a single byte through the socket.\n" +
				"@param value int Byte value to send (e.g. 10=newline, 2=STX, 3=ETX)\n" +
				"@param socket_name string Socket name\n" +
				"@returns bool True if sent successfully",
		},
		{
			Name:      "socket_send_int",
			Signature: "socket_send_int(value, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send a 32-bit integer in network byte order through the socket.\n" +
				"@param value int Integer value to send\n" +
				"@param socket_name string Socket name\n" +
				"@returns bool True if sent successfully",
		},
		{
			Name:      "socket_send_line",
			Signature: "socket_send_line(str, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send a string followed by a newline through the socket.\n" +
				"@param str string String to send\n" +
				"@param socket_name string Socket name\n" +
				"@returns bool True if sent successfully",
		},
		{
			Name:      "socket_send_string",
			Signature: "socket_send_string(str, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send a string through the socket (no newline appended).\n" +
				"@param str string String to send\n" +
				"@param socket_name string Socket name\n" +
				"@returns bool True if sent successfully",
		},
		{
			Name:      "socket_set_var",
			Signature: "socket_set_var(name, value, socket_name=\"socket_0\")",
			Module:    "interfaces",
			Doc: "Send \"SET <name> <value>\\n\" through the socket.\n" +
				"@param name string Variable name\n" +
				"@param value int Value to send\n" +
				"@param socket_name string Socket name",
		},
		// F/T sensor
		{
			Name:      "zero_ftsensor",
			Signature: "zero_ftsensor()",
			Module:    "interfaces",
			Doc:       "Zero the TCP force/torque sensor by subtracting the current measurement from future readings.",
		},
		// Operator input
		{
			Name:      "request_boolean_from_primary_client",
			Signature: "request_boolean_from_primary_client(message)",
			Module:    "interfaces",
			Doc: "Show a yes/no dialog on PolyScope and wait for operator response.\n" +
				"@param message string Message to display\n" +
				"@returns bool True for yes, False for no",
		},
		{
			Name:      "request_float_from_primary_client",
			Signature: "request_float_from_primary_client(message)",
			Module:    "interfaces",
			Doc: "Show a decimal number input dialog on PolyScope and wait for operator response.\n" +
				"@param message string Message to display\n" +
				"@returns float Value entered by the operator",
		},
		{
			Name:      "request_integer_from_primary_client",
			Signature: "request_integer_from_primary_client(message)",
			Module:    "interfaces",
			Doc: "Show an integer input dialog on PolyScope and wait for operator response.\n" +
				"@param message string Message to display\n" +
				"@returns int Value entered by the operator",
		},
		{
			Name:      "request_string_from_primary_client",
			Signature: "request_string_from_primary_client(message)",
			Module:    "interfaces",
			Doc: "Show a text input dialog on PolyScope and wait for operator response.\n" +
				"@param message string Message to display\n" +
				"@returns string Value entered by the operator",
		},
		// Threads (language-level but commonly used)
		{
			Name:      "run",
			Signature: "run(thread_fn)",
			Module:    "interfaces",
			Doc: "Start a new thread executing the given thread function.\n" +
				"@param thread_fn function Thread function to execute\n" +
				"@returns handle Thread handle for use with kill()",
		},
		{
			Name:      "kill",
			Signature: "kill(thread_handle)",
			Module:    "interfaces",
			Doc: "Terminate a running thread.\n" +
				"@param thread_handle handle Thread handle returned by run()",
		},
		{
			Name:      "enter_critical",
			Signature: "enter_critical()",
			Module:    "interfaces",
			Doc:       "Enter a critical section, blocking thread switching until exit_critical is called.",
		},
		{
			Name:      "exit_critical",
			Signature: "exit_critical()",
			Module:    "interfaces",
			Doc:       "Exit a critical section, re-enabling thread switching.",
		},
	}
}
