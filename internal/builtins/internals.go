package builtins

func internalsBuiltins() []Builtin {
	return []Builtin{
		{
			Name:      "force",
			Signature: "force()",
			Module:    "internals",
			Doc: "Returns the norm of the external TCP force (Fx, Fy, Fz) computed from get_tcp_force.\n" +
				"@returns float Force in Newtons",
		},
		{
			Name:      "estimate_payload",
			Signature: "estimate_payload(poses, wrenches)",
			Module:    "internals",
			Doc: "Estimate payload mass and center of gravity from a set of poses and corresponding wrenches.\n" +
				"Requires at least four poses with varied orientations.\n" +
				"@param poses list List of TCP poses recorded at different orientations\n" +
				"@param wrenches list Corresponding wrenches measured at each pose\n" +
				"@returns struct Struct with mass (kg) and cog (3d vector in tool frame [m])",
		},
		{
			Name:      "get_actual_joint_positions",
			Signature: "get_actual_joint_positions()",
			Module:    "internals",
			Doc: "Returns the actual angular positions of all joints.\n" +
				"@returns list Joint positions [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in rad",
		},
		{
			Name:      "get_actual_joint_positions_history",
			Signature: "get_actual_joint_positions_history(steps=0)",
			Module:    "internals",
			Doc: "Returns the actual joint positions from a number of controller time steps ago.\n" +
				"@param steps int Number of time steps back (0 = current)\n" +
				"@returns list Joint positions [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in rad",
		},
		{
			Name:      "get_actual_joint_speeds",
			Signature: "get_actual_joint_speeds()",
			Module:    "internals",
			Doc: "Returns the actual angular velocities of all joints.\n" +
				"@returns list Joint speeds [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in rad/s",
		},
		{
			Name:      "get_actual_tcp_pose",
			Signature: "get_actual_tcp_pose()",
			Module:    "internals",
			Doc: "Returns the current measured tool pose based on encoder readings.\n" +
				"@returns pose TCP vector [X, Y, Z, Rx, Ry, Rz]",
		},
		{
			Name:      "get_actual_tcp_speed",
			Signature: "get_actual_tcp_speed()",
			Module:    "internals",
			Doc: "Returns the current measured TCP speed.\n" +
				"@returns list Speed vector [X, Y, Z, Rx, Ry, Rz] in m/s and rad/s",
		},
		{
			Name:      "get_actual_tool_flange_pose",
			Signature: "get_actual_tool_flange_pose()",
			Module:    "internals",
			Doc: "Returns the tool flange pose without TCP offset, based on encoder readings.\n" +
				"@returns pose Flange vector [X, Y, Z, Rx, Ry, Rz]",
		},
		{
			Name:      "get_controller_temp",
			Signature: "get_controller_temp()",
			Module:    "internals",
			Doc: "Returns the control box temperature.\n" +
				"@returns float Temperature in degrees Celsius",
		},
		{
			Name:      "get_forward_kin",
			Signature: "get_forward_kin(q=\"current_joint_positions\", tcp=\"active_tcp\")",
			Module:    "internals",
			Doc: "Compute forward kinematics (joint space to tool space) using calibrated kinematics.\n" +
				"@param q list Joint position vector (optional, defaults to current)\n" +
				"@param tcp pose TCP offset pose (optional, defaults to active)\n" +
				"@returns pose Tool pose",
		},
		{
			Name:      "get_inverse_kin",
			Signature: "get_inverse_kin(x, qnear, maxPositionError=1e-10, maxOrientationError=1e-10, tcp=\"active_tcp\")",
			Module:    "internals",
			Doc: "Compute inverse kinematics (tool space to joint space).\n" +
				"Returns the solution closest to qnear, or to current joint positions if qnear is omitted.\n" +
				"@param x pose Target tool pose\n" +
				"@param qnear list Preferred joint positions for nearest solution (optional)\n" +
				"@param maxPositionError float Maximum allowed position error [m]\n" +
				"@param maxOrientationError float Maximum allowed orientation error [rad]\n" +
				"@param tcp pose TCP offset (optional, defaults to active)\n" +
				"@returns list Joint positions",
		},
		{
			Name:      "get_inverse_kin_has_solution",
			Signature: "get_inverse_kin_has_solution(pose, qnear, maxPositionError=1e-10, maxOrientationError=1e-10, tcp=\"active_tcp\")",
			Module:    "internals",
			Doc: "Check whether get_inverse_kin has a valid solution for the given pose.\n" +
				"@param pose pose Target tool pose\n" +
				"@param qnear list Preferred joint positions (optional)\n" +
				"@param maxPositionError float Maximum allowed position error [m]\n" +
				"@param maxOrientationError float Maximum allowed orientation error [rad]\n" +
				"@param tcp pose TCP offset (optional)\n" +
				"@returns bool True if a solution exists",
		},
		{
			Name:      "get_joint_temp",
			Signature: "get_joint_temp(j)",
			Module:    "internals",
			Doc: "Returns the temperature of the specified joint housing.\n" +
				"@param j int Joint number (0=base, 5=wrist3)\n" +
				"@returns float Temperature in degrees Celsius",
		},
		{
			Name:      "get_joint_torques",
			Signature: "get_joint_torques()",
			Module:    "internals",
			Doc: "Returns the external torques on all joints, corrected for gravity and friction.\n" +
				"@returns list Torques [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in Nm",
		},
		{
			Name:      "get_steptime",
			Signature: "get_steptime()",
			Module:    "internals",
			Doc: "Returns the duration of the robot controller time step.\n" +
				"@returns float Step duration in seconds",
		},
		{
			Name:      "get_target_joint_positions",
			Signature: "get_target_joint_positions()",
			Module:    "internals",
			Doc: "Returns the desired (target) angular positions of all joints.\n" +
				"@returns list Joint positions [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in rad",
		},
		{
			Name:      "get_target_joint_speeds",
			Signature: "get_target_joint_speeds()",
			Module:    "internals",
			Doc: "Returns the desired (target) angular velocities of all joints.\n" +
				"@returns list Joint speeds [Base, Shoulder, Elbow, Wrist1, Wrist2, Wrist3] in rad/s",
		},
		{
			Name:      "get_target_payload",
			Signature: "get_target_payload()",
			Module:    "internals",
			Doc: "Returns the weight of the active payload.\n" +
				"@returns float Payload mass in kilograms",
		},
		{
			Name:      "get_target_payload_cog",
			Signature: "get_target_payload_cog()",
			Module:    "internals",
			Doc: "Returns the center of gravity of the active payload relative to the tool flange.\n" +
				"@returns list CoG coordinates [x, y, z] in meters",
		},
		{
			Name:      "get_target_payload_inertia",
			Signature: "get_target_payload_inertia()",
			Module:    "internals",
			Doc: "Returns the inertia matrix of the active payload in tool flange coordinates.\n" +
				"@returns list Inertia [Ixx, Iyy, Izz, Ixy, Ixz, Iyz] in kg*m^2",
		},
		{
			Name:      "get_target_tcp_pose",
			Signature: "get_target_tcp_pose()",
			Module:    "internals",
			Doc: "Returns the current target tool pose based on target joint positions.\n" +
				"@returns pose TCP vector [X, Y, Z, Rx, Ry, Rz]",
		},
		{
			Name:      "get_target_tcp_speed",
			Signature: "get_target_tcp_speed()",
			Module:    "internals",
			Doc: "Returns the current desired TCP speed.\n" +
				"@returns pose Speed vector [X, Y, Z, Rx, Ry, Rz]",
		},
		{
			Name:      "get_target_waypoint",
			Signature: "get_target_waypoint()",
			Module:    "internals",
			Doc: "Returns the target waypoint of the active move command.\n" +
				"Unlike get_target_tcp_pose, this returns the same pose throughout the entire motion.\n" +
				"@returns pose Waypoint TCP vector [X, Y, Z, Rx, Ry, Rz]",
		},
		{
			Name:      "get_tcp_force",
			Signature: "get_tcp_force()",
			Module:    "internals",
			Doc: "Returns the wrench (force and torque) at the tool flange in base frame orientation.\n" +
				"@returns pose [Fx(N), Fy(N), Fz(N), TRx(Nm), TRy(Nm), TRz(Nm)]",
		},
		{
			Name:      "get_tcp_offset",
			Signature: "get_tcp_offset()",
			Module:    "internals",
			Doc: "Returns the active TCP offset from the tool flange to the TCP.\n" +
				"@returns pose TCP offset pose",
		},
		{
			Name:      "get_tool_accelerometer_reading",
			Signature: "get_tool_accelerometer_reading()",
			Module:    "internals",
			Doc: "Returns the current tool accelerometer reading aligned with tool coordinates.\n" +
				"@returns list Acceleration [x, y, z] in m/s^2",
		},
		{
			Name:      "get_tool_current",
			Signature: "get_tool_current()",
			Module:    "internals",
			Doc: "Returns the tool current consumption.\n" +
				"@returns float Current in amperes",
		},
		{
			Name:      "high_holding_torque_disable",
			Signature: "high_holding_torque_disable()",
			Module:    "internals",
			Doc:       "Disable automatic high holding torque when the robot is stationary. Resets on controller restart.",
		},
		{
			Name:      "high_holding_torque_enable",
			Signature: "high_holding_torque_enable()",
			Module:    "internals",
			Doc:       "Re-enable high holding torque when stationary (reverses high_holding_torque_disable).",
		},
		{
			Name:      "is_steady",
			Signature: "is_steady()",
			Module:    "internals",
			Doc: "Returns True when the robot has been at rest with zero target velocity for 500ms.\n" +
				"@returns bool True when steady",
		},
		{
			Name:      "is_within_safety_limits",
			Signature: "is_within_safety_limits(pose, qnear)",
			Module:    "internals",
			Doc: "Deprecated: Check if a pose is reachable within current safety limits. Use get_inverse_kin_has_solution instead.\n" +
				"@param pose pose Target pose or joint positions\n" +
				"@param qnear list Joint positions hint (optional, only for pose targets)\n" +
				"@returns bool True if within limits",
		},
		{
			Name:      "popup",
			Signature: "popup(s, title=\"Popup\", warning=False, error=False, blocking=False)",
			Module:    "internals",
			Doc: "Display a popup message on the teach pendant.\n" +
				"@param s string Message text\n" +
				"@param title string Popup title\n" +
				"@param warning bool Show as warning\n" +
				"@param error bool Show as error\n" +
				"@param blocking bool If True, program pauses until user presses continue",
		},
		{
			Name:      "powerdown",
			Signature: "powerdown()",
			Module:    "internals",
			Doc:       "Shut down the robot and power off the controller.",
		},
		{
			Name:      "protective_stop",
			Signature: "protective_stop()",
			Module:    "internals",
			Doc:       "Trigger a protective stop, pausing the program and halting motion.",
		},
		{
			Name:      "set_base_acceleration",
			Signature: "set_base_acceleration(a)",
			Module:    "internals",
			Doc: "Specify the linear acceleration of the robot base (for robots on moving platforms).\n" +
				"@param a list Base acceleration [x, y, z] in m/s^2",
		},
		{
			Name:      "set_baselight_off",
			Signature: "set_baselight_off()",
			Module:    "internals",
			Doc:       "Turn off the baselight ring. Only applies to UR20/UR30.",
		},
		{
			Name:      "set_baselight_iec",
			Signature: "set_baselight_iec()",
			Module:    "internals",
			Doc:       "Set baselight to IEC 60204-1 standard mode indicating safety state. Only applies to UR20/UR30.",
		},
		{
			Name:      "set_baselight_solid",
			Signature: "set_baselight_solid(r, g, b)",
			Module:    "internals",
			Doc: "Set a solid color on the baselight ring. Only applies to UR20/UR30.\n" +
				"@param r int Red component [0-255]\n" +
				"@param g int Green component [0-255]\n" +
				"@param b int Blue component [0-255]",
		},
		{
			Name:      "set_gravity",
			Signature: "set_gravity(d)",
			Module:    "internals",
			Doc: "Set the direction of gravity relative to the robot base.\n" +
				"@param d list 3D gravity vector [x, y, z] in m/s^2 (e.g. [0, 0, 9.82] for upright mounting)",
		},
		{
			Name:      "set_payload",
			Signature: "set_payload(m, cog)",
			Module:    "internals",
			Doc: "Deprecated: Set payload mass and center of gravity. Use set_target_payload instead.\n" +
				"Resets the payload inertia matrix. F/T sensor is automatically zeroed.\n" +
				"@param m float Mass in kilograms\n" +
				"@param cog list Center of gravity offset [x, y, z] in meters from tool mount",
		},
		{
			Name:      "set_payload_cog",
			Signature: "set_payload_cog(CoG)",
			Module:    "internals",
			Doc: "Deprecated: Set the payload center of gravity. Use set_target_payload instead.\n" +
				"@param CoG list Center of gravity [x, y, z] in meters from tool mount",
		},
		{
			Name:      "set_payload_mass",
			Signature: "set_payload_mass(m)",
			Module:    "internals",
			Doc: "Deprecated: Set payload mass, keeping CoG unchanged. Use set_target_payload instead.\n" +
				"@param m float Mass in kilograms",
		},
		{
			Name:      "set_target_payload",
			Signature: "set_target_payload(m, cog, inertia=[0,0,0,0,0,0], transition_time=0)",
			Module:    "internals",
			Doc: "Set payload mass, center of gravity, inertia matrix, and transition time.\n" +
				"Call when the robot picks up or puts down a workpiece. F/T sensor is automatically zeroed.\n" +
				"@param m float Mass in kilograms\n" +
				"@param cog list Center of gravity [CoGx, CoGy, CoGz] in meters from tool mount\n" +
				"@param inertia list Inertia matrix [Ixx, Iyy, Izz, Ixy, Ixz, Iyz] in kg*m^2\n" +
				"@param transition_time float Duration of the payload change in seconds (0 = instant)",
		},
		{
			Name:      "set_tcp",
			Signature: "set_tcp(pose, tcp_name=\"\")",
			Module:    "internals",
			Doc: "Set the active tool center point offset from the tool flange.\n" +
				"@param pose pose TCP offset transformation\n" +
				"@param tcp_name string Name to assign to this TCP configuration",
		},
		{
			Name:      "sleep",
			Signature: "sleep(t)",
			Module:    "internals",
			Doc: "Pause program execution for the specified duration.\n" +
				"@param t float Time in seconds",
		},
		{
			Name:      "str_at",
			Signature: "str_at(src, index)",
			Module:    "internals",
			Doc: "Returns the byte at the given index in a string.\n" +
				"@param src string Source string\n" +
				"@param index int Zero-based position\n" +
				"@returns string Single-byte string at the given position",
		},
		{
			Name:      "str_cat",
			Signature: "str_cat(op1, op2)",
			Module:    "internals",
			Doc: "Concatenate two values as strings. Accepts strings, bools, numbers, poses, and lists.\n" +
				"@param op1 any First operand\n" +
				"@param op2 any Second operand\n" +
				"@returns string Concatenated result (max 1023 characters)",
		},
		{
			Name:      "str_empty",
			Signature: "str_empty(str)",
			Module:    "internals",
			Doc: "Check whether a string is empty.\n" +
				"@param str string Source string\n" +
				"@returns bool True if the string has zero length",
		},
		{
			Name:      "str_find",
			Signature: "str_find(src, target, start_from=0)",
			Module:    "internals",
			Doc: "Find the first occurrence of a substring.\n" +
				"@param src string Source string to search in\n" +
				"@param target string Substring to find\n" +
				"@param start_from int Starting position (default 0)\n" +
				"@returns int Index of first match, or -1 if not found",
		},
		{
			Name:      "str_len",
			Signature: "str_len(str)",
			Module:    "internals",
			Doc: "Returns the number of bytes in a string.\n" +
				"@param str string Source string\n" +
				"@returns int Byte count",
		},
		{
			Name:      "str_sub",
			Signature: "str_sub(src, index, len)",
			Module:    "internals",
			Doc: "Returns a substring starting at the given index.\n" +
				"@param src string Source string\n" +
				"@param index int Starting byte position (zero-based)\n" +
				"@param len int Maximum length of substring (optional, defaults to end of string)\n" +
				"@returns string The extracted substring",
		},
		{
			Name:      "sync",
			Signature: "sync()",
			Module:    "internals",
			Doc:       "Uses up the remaining physical time in the current frame, allowing other threads to run.",
		},
		{
			Name:      "textmsg",
			Signature: "textmsg(s1, s2=\"\")",
			Module:    "internals",
			Doc: "Send a text message to the PolyScope log tab.\n" +
				"@param s1 any Main message (string, int, bool, pose, etc.)\n" +
				"@param s2 any Secondary message, concatenated after s1",
		},
		{
			Name:      "to_num",
			Signature: "to_num(str)",
			Module:    "internals",
			Doc: "Convert a string to a number. Returns int or float depending on decimal point presence.\n" +
				"Supports decimal, hex (0x), infinity, and NaN.\n" +
				"@param str string String to convert\n" +
				"@returns number Integer or float value",
		},
		{
			Name:      "to_str",
			Signature: "to_str(val)",
			Module:    "internals",
			Doc: "Convert a value to its string representation. Floats use 6 decimal places.\n" +
				"@param val any Value to convert (bool, int, float, pose, or list)\n" +
				"@returns string String representation (max 1023 characters)",
		},
		{
			Name:      "tool_contact",
			Signature: "tool_contact(direction)",
			Module:    "internals",
			Doc: "Detect contact between the tool and an object.\n" +
				"@param direction list 6d vector giving the detection direction in base coordinates (all zeros = any direction)\n" +
				"@returns int Number of time steps back to contact start (0 = no contact)",
		},
	}
}
