package builtins

func ioconfigurationBuiltins() []Builtin {
	return []Builtin{
		{
			Name:      "modbus_set_runstate_dependent_choice",
			Signature: "modbus_set_runstate_dependent_choice(signal_name, runstate_choice)",
			Module:    "ioconfiguration",
			Doc: "Set Modbus output signal behavior based on program state.\n" +
				"@param signal_name string Name of a digital or register output signal\n" +
				"@param runstate_choice int Behavior (0=preserve, 1=low when stopped, 2=high when stopped, 3=high when running, 4=low on fault, 5=high from start/low on fault)",
		},
		{
			Name:      "set_analog_outputdomain",
			Signature: "set_analog_outputdomain(port, domain)",
			Module:    "ioconfiguration",
			Doc: "Set the domain of an analog output.\n" +
				"@param port int Analog output port number\n" +
				"@param domain int Output domain (0=4-20mA, 1=0-10V)",
		},
		{
			Name:      "set_configurable_digital_input_action",
			Signature: "set_configurable_digital_input_action(port, action)",
			Module:    "ioconfiguration",
			Doc: "Set a configurable digital input to a default or freedrive action.\n" +
				"@param port int Input port number\n" +
				"@param action string Action type: \"default\" or \"freedrive\"",
		},
		{
			Name:      "set_gp_boolean_input_action",
			Signature: "set_gp_boolean_input_action(port, action)",
			Module:    "ioconfiguration",
			Doc: "Set a general purpose boolean input to a default or freedrive action.\n" +
				"@param port int Input port number [0-127]\n" +
				"@param action string Action type: \"default\" or \"freedrive\"",
		},
		{
			Name:      "set_input_actions_to_default",
			Signature: "set_input_actions_to_default()",
			Module:    "ioconfiguration",
			Doc:       "Reset all standard, configurable, tool, and gp_boolean input actions to default.",
		},
		{
			Name:      "set_runstate_configurable_digital_output_to_value",
			Signature: "set_runstate_configurable_digital_output_to_value(outputId, state)",
			Module:    "ioconfiguration",
			Doc: "Set configurable digital output behavior based on program state.\n" +
				"@param outputId int Output number [0-7]\n" +
				"@param state int Behavior (0=preserve, 1=low when stopped, 2=high when stopped, 3=high when running, 4=low on fault, 5=high from start/low on fault)",
		},
		{
			Name:      "set_runstate_gp_boolean_output_to_value",
			Signature: "set_runstate_gp_boolean_output_to_value(outputId, state)",
			Module:    "ioconfiguration",
			Doc: "Set general purpose boolean output behavior based on program state.\n" +
				"@param outputId int Output number [0-127]\n" +
				"@param state int Behavior (0=preserve, 1=false when stopped, 2=true when stopped, 3=true when running, 4=false on fault, 5=true from start/false on fault)",
		},
		{
			Name:      "set_runstate_standard_analog_output_to_value",
			Signature: "set_runstate_standard_analog_output_to_value(outputId, state)",
			Module:    "ioconfiguration",
			Doc: "Set standard analog output behavior based on program state.\n" +
				"@param outputId int Output number [0-1]\n" +
				"@param state int Behavior (0=preserve, 1=low when stopped, 2=high when stopped, 3=high when running, 4=low on fault, 5=high from start/low on fault)",
		},
		{
			Name:      "set_runstate_standard_digital_output_to_value",
			Signature: "set_runstate_standard_digital_output_to_value(outputId, state)",
			Module:    "ioconfiguration",
			Doc: "Set standard digital output behavior based on program state.\n" +
				"@param outputId int Output number [0-7]\n" +
				"@param state int Behavior (0=preserve, 1=low when stopped, 2=high when stopped, 3=high when running, 4=low on fault, 5=high from start/low on fault)",
		},
		{
			Name:      "set_runstate_tool_digital_output_to_value",
			Signature: "set_runstate_tool_digital_output_to_value(outputId, state)",
			Module:    "ioconfiguration",
			Doc: "Set tool digital output behavior based on program state.\n" +
				"@param outputId int Output number [0-1]\n" +
				"@param state int Behavior (0=preserve, 1=low when stopped, 2=high when stopped, 3=high when running, 4=low on fault, 5=high from start/low on fault)",
		},
		{
			Name:      "set_standard_analog_input_domain",
			Signature: "set_standard_analog_input_domain(port, domain)",
			Module:    "ioconfiguration",
			Doc: "Set the domain of a standard analog input in the controller box.\n" +
				"@param port int Analog input port (0 or 1)\n" +
				"@param domain int Input domain (0=4-20mA, 1=0-10V)",
		},
		{
			Name:      "set_standard_digital_input_action",
			Signature: "set_standard_digital_input_action(port, action)",
			Module:    "ioconfiguration",
			Doc: "Set a standard digital input to a default or freedrive action.\n" +
				"@param port int Input port number\n" +
				"@param action string Action type: \"default\" or \"freedrive\"",
		},
		{
			Name:      "set_tool_analog_input_domain",
			Signature: "set_tool_analog_input_domain(port, domain)",
			Module:    "ioconfiguration",
			Doc: "Set the domain of a tool analog input.\n" +
				"@param port int Analog input port (0 or 1)\n" +
				"@param domain int Input domain (0=4-20mA, 1=0-10V)",
		},
		{
			Name:      "set_tool_digital_input_action",
			Signature: "set_tool_digital_input_action(port, action)",
			Module:    "ioconfiguration",
			Doc: "Set a tool digital input to a default or freedrive action.\n" +
				"@param port int Input port number\n" +
				"@param action string Action type: \"default\" or \"freedrive\"",
		},
	}
}
