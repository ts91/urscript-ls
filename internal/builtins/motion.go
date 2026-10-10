package builtins

func motionBuiltins() []Builtin {
	return []Builtin{
		{
			Name:      "conveyor_pulse_decode",
			Signature: "conveyor_pulse_decode(type, A, B)",
			Module:    "motion",
			Doc: "Deprecated: Configure conveyor encoder pulse decoding on digital inputs.\n" +
				"@param type int Encoder type (0=none, 1=quadrature, 2=rising+falling, 3=rising, 4=falling)\n" +
				"@param A int Encoder input A, digital input number [0-3]\n" +
				"@param B int Encoder input B, digital input number [0-3]",
		},
		{
			Name:      "encoder_enable_pulse_decode",
			Signature: "encoder_enable_pulse_decode(encoder_index, decoder_type, A, B)",
			Module:    "motion",
			Doc: "Configure an encoder connected to the pulse decoder.\n" +
				"@param encoder_index int Index of the encoder, must be 0 or 1\n" +
				"@param decoder_type int Decoder type (0=none, 1=quadrature, 2=rising+falling, 3=rising, 4=falling)\n" +
				"@param A int Encoder input A pin [8-11]\n" +
				"@param B int Encoder input B pin [8-11]",
		},
		{
			Name:      "encoder_enable_set_tick_count",
			Signature: "encoder_enable_set_tick_count(encoder_index, range_id)",
			Module:    "motion",
			Doc: "Configure an encoder to receive tick counts via encoder_set_tick_count.\n" +
				"@param encoder_index int Index of the encoder, must be 0 or 1\n" +
				"@param range_id int Encoder range (0=32-bit signed, 1=8-bit unsigned, 2=16-bit unsigned, 3=24-bit unsigned, 4=32-bit unsigned)",
		},
		{
			Name:      "encoder_get_tick_count",
			Signature: "encoder_get_tick_count(encoder_index, opt=\"\")",
			Module:    "motion",
			Doc: "Read the filtered tick count of an encoder.\n" +
				"@param encoder_index int Index of the encoder, must be 0 or 1\n" +
				"@param opt string Option flag, use \"raw\" for unfiltered integer value\n" +
				"@returns float Filtered tick count, or int if opt=\"raw\"",
		},
		{
			Name:      "encoder_set_tick_count",
			Signature: "encoder_set_tick_count(encoder_index, count)",
			Module:    "motion",
			Doc: "Set the tick count of an encoder. Useful for absolute encoders.\n" +
				"@param encoder_index int Index of the encoder, must be 0 or 1\n" +
				"@param count int Tick count value within the encoder range",
		},
		{
			Name:      "encoder_unwind_delta_tick_count",
			Signature: "encoder_unwind_delta_tick_count(encoder_index, delta_tick_count)",
			Module:    "motion",
			Doc: "Compensate for encoder wrapping when computing tick count deltas.\n" +
				"If the delta exceeds half the encoder range, wrapping is assumed and corrected.\n" +
				"@param encoder_index int Index of the encoder, must be 0 or 1\n" +
				"@param delta_tick_count float Difference between two tick counts\n" +
				"@returns float Unwound delta tick count",
		},
		{
			Name:      "end_force_mode",
			Signature: "end_force_mode()",
			Module:    "motion",
			Doc:       "Exit force mode and return to normal position control.",
		},
		{
			Name:      "end_freedrive_mode",
			Signature: "end_freedrive_mode()",
			Module:    "motion",
			Doc:       "Exit freedrive mode and return to normal position control.",
		},
		{
			Name:      "end_screw_driving",
			Signature: "end_screw_driving()",
			Module:    "motion",
			Doc:       "Exit screw driving mode and return to normal operation.",
		},
		{
			Name:      "end_teach_mode",
			Signature: "end_teach_mode()",
			Module:    "motion",
			Doc:       "Deprecated: Exit teach mode. Use end_freedrive_mode instead.",
		},
		{
			Name:      "force_mode",
			Signature: "force_mode(task_frame, selection_vector, wrench, type, limits)",
			Module:    "motion",
			Doc: "Set robot to force mode where it complies in specified axes.\n" +
				"@param task_frame pose Force frame relative to the base frame\n" +
				"@param selection_vector list 6d vector of 0s and 1s selecting compliant axes\n" +
				"@param wrench list Target forces/torques the robot applies to its environment\n" +
				"@param type int Force frame interpretation (1=TCP-to-origin aligned, 2=unchanged, 3=velocity aligned)\n" +
				"@param limits list 6d vector of max TCP speed (compliant) or max deviation (non-compliant)",
		},
		{
			Name:      "force_mode_set_damping",
			Signature: "force_mode_set_damping(damping)",
			Module:    "motion",
			Doc: "Set the damping parameter in force mode.\n" +
				"@param damping float Value between 0 (no damping) and 1 (full damping), default 0.005",
		},
		{
			Name:      "force_mode_set_gain_scaling",
			Signature: "force_mode_set_gain_scaling(scaling)",
			Module:    "motion",
			Doc: "Scale the gain in force mode. Values above 1 may cause instability.\n" +
				"@param scaling float Gain scaling between 0 and 2, default 1",
		},
		{
			Name:      "freedrive_mode",
			Signature: "freedrive_mode(freeAxes=[1,1,1,1,1,1], feature=p[0,0,0,0,0,0])",
			Module:    "motion",
			Doc: "Set robot in freedrive mode, allowing it to be moved by hand.\n" +
				"Use constrained freedrive by specifying which axes allow movement.\n" +
				"@param freeAxes list 6d vector of 0s and 1s selecting free axes [x,y,z,rx,ry,rz]\n" +
				"@param feature pose Reference frame for the free axes, relative to base. Use \"base\" or \"tool\" for predefined frames",
		},
		{
			Name:      "get_conveyor_tick_count",
			Signature: "get_conveyor_tick_count()",
			Module:    "motion",
			Doc: "Deprecated: Read the interpolated conveyor encoder tick count. Use encoder_get_tick_count instead.\n" +
				"@returns float Conveyor encoder tick count",
		},
		{
			Name:      "get_freedrive_status",
			Signature: "get_freedrive_status()",
			Module:    "motion",
			Doc: "Query the freedrive status relative to singularity proximity.\n" +
				"@returns int 0=normal, 1=near singularity, 2=too close to singularity",
		},
		{
			Name:      "get_target_tcp_pose_along_path",
			Signature: "get_target_tcp_pose_along_path()",
			Module:    "motion",
			Doc: "Query the target TCP pose from the active trajectory without any path offset applied.\n" +
				"@returns pose Target TCP pose",
		},
		{
			Name:      "get_target_tcp_speed_along_path",
			Signature: "get_target_tcp_speed_along_path()",
			Module:    "motion",
			Doc: "Query the target TCP speed from the active trajectory without any path offset applied.\n" +
				"@returns list Target TCP speed vector",
		},
		{
			Name:      "movec",
			Signature: "movec(pose_via, pose_to, a=1.2, v=0.25, r=0, mode=0)",
			Module:    "motion",
			Doc: "Move circular: the TCP follows a circular arc through pose_via to pose_to.\n" +
				"@param pose_via pose Via point on the arc (only position is used)\n" +
				"@param pose_to pose Target pose (only position used in fixed mode)\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param r float Blend radius at target [m]\n" +
				"@param mode int Orientation mode (0=unconstrained interpolation, 1=fixed relative to arc tangent)",
		},
		{
			Name:      "movej",
			Signature: "movej(q, a=1.4, v=1.05, t=0, r=0)",
			Module:    "motion",
			Doc: "Move to position (linear in joint-space).\n" +
				"Robot must be at standstill or coming from a blended move. The t parameter overrides a and v when set.\n" +
				"@param q list Joint positions [rad] or a pose (inverse kinematics applied)\n" +
				"@param a float Joint acceleration of leading axis [rad/s^2]\n" +
				"@param v float Joint speed of leading axis [rad/s]\n" +
				"@param t float Duration of the move [s], 0 means not specified\n" +
				"@param r float Blend radius [m]",
		},
		{
			Name:      "movel",
			Signature: "movel(pose, a=1.2, v=0.25, t=0, r=0)",
			Module:    "motion",
			Doc: "Move to position (linear in tool-space).\n" +
				"@param pose pose Target pose (or joint positions, forward kinematics applied)\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param t float Duration of the move [s], 0 means not specified\n" +
				"@param r float Blend radius [m]",
		},
		{
			Name:      "movep",
			Signature: "movep(pose, a=1.2, v=0.25, r=0)",
			Module:    "motion",
			Doc: "Move process: blend circular and move linear in tool-space at constant speed.\n" +
				"@param pose pose Target pose (or joint positions, forward kinematics applied)\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param r float Blend radius [m]",
		},
		{
			Name:      "path_offset_disable",
			Signature: "path_offset_disable(a=20)",
			Module:    "motion",
			Doc: "Disable path offsetting and decelerate all joints to zero speed.\n" +
				"@param a float Joint deceleration [rad/s^2]",
		},
		{
			Name:      "path_offset_enable",
			Signature: "path_offset_enable()",
			Module:    "motion",
			Doc: "Enable path offsetting to superimpose a Cartesian offset onto robot motion.\n" +
				"Useful for weaving patterns in welding or compensating for base movement.",
		},
		{
			Name:      "path_offset_get",
			Signature: "path_offset_get(type)",
			Module:    "motion",
			Doc: "Query the currently applied path offset.\n" +
				"@param type int Frame of reference (see path_offset_set)\n" +
				"@returns pose Translational and rotational offset [m, rad]",
		},
		{
			Name:      "path_offset_set",
			Signature: "path_offset_set(offset, type)",
			Module:    "motion",
			Doc: "Specify the Cartesian path offset to apply. Call path_offset_enable first.\n" +
				"@param offset pose Translational and rotational offset\n" +
				"@param type int Reference frame (1=BASE, 2=TCP, 3=MOTION along trajectory tangent, 4=WORLD for mobile bases)",
		},
		{
			Name:      "path_offset_set_alpha_filter",
			Signature: "path_offset_set_alpha_filter(alpha)",
			Module:    "motion",
			Doc: "Set the alpha filter coefficient for path offset smoothing (EWMA).\n" +
				"@param alpha float Filter coefficient between 0 and 1. 1 = no filtering, ~0.1 is typical for welding",
		},
		{
			Name:      "path_offset_set_max_offset",
			Signature: "path_offset_set_max_offset(transLimit, rotLimit)",
			Module:    "motion",
			Doc: "Set limits for the maximum allowed path offset.\n" +
				"@param transLimit float Maximum translational offset along any axis [m]\n" +
				"@param rotLimit float Maximum rotational offset around any axis [rad]",
		},
		{
			Name:      "pause_on_error_code",
			Signature: "pause_on_error_code(code, argument)",
			Module:    "motion",
			Doc: "Pause the program if the specified error code occurs during execution.\n" +
				"@param code int Error code to watch for\n" +
				"@param argument int Error argument to match (omit to match any argument)",
		},
		{
			Name:      "position_deviation_warning",
			Signature: "position_deviation_warning(enabled, threshold=0.8)",
			Module:    "motion",
			Doc: "Enable or disable logging of position deviation warnings.\n" +
				"@param enabled bool Enable or disable warning messages\n" +
				"@param threshold float Deviation level in [0,1] that triggers a warning (1 = protective stop level)",
		},
		{
			Name:      "reset_revolution_counter",
			Signature: "reset_revolution_counter(qNear=[0.0,0.0,0.0,0.0,0.0,0.0])",
			Module:    "motion",
			Doc: "Reset the revolution counter for joints with unlimited safety limits.\n" +
				"@param qNear list Joint vector to reset near (uses actual revolutions if omitted)",
		},
		{
			Name:      "screw_driving",
			Signature: "screw_driving(f, v_limit)",
			Module:    "motion",
			Doc: "Enter screw driving mode. The robot exerts force along TCP Z-axis at limited speed.\n" +
				"Zero the F/T sensor before calling. Call end_screw_driving when done.\n" +
				"@param f float Force along TCP Z-axis [N]\n" +
				"@param v_limit float Maximum TCP velocity along Z [m/s]",
		},
		{
			Name:      "servoc",
			Signature: "servoc(pose, a=1.2, v=0.25, r=0)",
			Module:    "motion",
			Doc: "Servo to position (circular in tool-space) at constant tool speed.\n" +
				"@param pose pose Target pose (or joint positions, forward kinematics applied)\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param r float Blend radius [m]",
		},
		{
			Name:      "servoj",
			Signature: "servoj(q, a, v, t=0.002, lookahead_time=0.1, gain=300)",
			Module:    "motion",
			Doc: "Servo to position in joint-space for online real-time control.\n" +
				"Call with a new setpoint each time step. Combine with get_inverse_kin for Cartesian servoing.\n" +
				"@param q list Joint positions [rad]\n" +
				"@param a float Not used in current version\n" +
				"@param v float Not used in current version\n" +
				"@param t float Time for the command to control the robot [s]\n" +
				"@param lookahead_time float Smoothing time [s], range [0.03, 0.2]\n" +
				"@param gain float Proportional gain for position tracking, range [100, 2000]",
		},
		{
			Name:      "set_conveyor_tick_count",
			Signature: "set_conveyor_tick_count(tick_count, absolute_encoder_resolution=0)",
			Module:    "motion",
			Doc: "Deprecated: Set the conveyor encoder tick count. Use encoder_set_tick_count instead.\n" +
				"@param tick_count int Encoder tick count\n" +
				"@param absolute_encoder_resolution int Encoder range (0=32-bit signed, 1=8-bit, 2=16-bit, 3=24-bit, 4=32-bit unsigned)",
		},
		{
			Name:      "set_pos",
			Signature: "set_pos(q)",
			Module:    "motion",
			Doc: "Set joint positions of the simulated robot.\n" +
				"@param q list Joint positions [rad]",
		},
		{
			Name:      "set_safety_mode_transition_hardness",
			Signature: "set_safety_mode_transition_hardness(type)",
			Module:    "motion",
			Doc: "Set the transition hardness between normal mode, reduced mode, and safeguard stop.\n" +
				"@param type int 0=hard transition (max torque), 1=soft transition",
		},
		{
			Name:      "speedj",
			Signature: "speedj(qd, a, t)",
			Module:    "motion",
			Doc: "Accelerate linearly in joint-space and continue at constant joint speed.\n" +
				"Returns when target speed is reached, or after time t if specified.\n" +
				"@param qd list Joint speeds [rad/s]\n" +
				"@param a float Joint acceleration of leading axis [rad/s^2]\n" +
				"@param t float Time before return [s] (optional)",
		},
		{
			Name:      "speedl",
			Signature: "speedl(xd, a, t, aRot=a)",
			Module:    "motion",
			Doc: "Accelerate linearly in Cartesian space and continue at constant tool speed.\n" +
				"Returns when target speed is reached, or after time t if specified.\n" +
				"@param xd list Tool speed [m/s] as spatial vector\n" +
				"@param a float Tool positional acceleration [m/s^2]\n" +
				"@param t float Time before return [s] (optional)\n" +
				"@param aRot float Tool rotational acceleration [rad/s^2] (defaults to a)",
		},
		{
			Name:      "stop_conveyor_tracking",
			Signature: "stop_conveyor_tracking(a=20)",
			Module:    "motion",
			Doc: "Stop conveyor tracking and decelerate all joints to zero speed.\n" +
				"@param a float Joint deceleration [rad/s^2]",
		},
		{
			Name:      "stopj",
			Signature: "stopj(a)",
			Module:    "motion",
			Doc: "Decelerate joint speeds to zero (linear in joint-space).\n" +
				"@param a float Joint deceleration of leading axis [rad/s^2]",
		},
		{
			Name:      "stopl",
			Signature: "stopl(a, aRot=a)",
			Module:    "motion",
			Doc: "Decelerate tool speed to zero (linear in tool-space).\n" +
				"@param a float Tool deceleration [m/s^2]\n" +
				"@param aRot float Tool rotational deceleration [rad/s^2] (defaults to a)",
		},
		{
			Name:      "teach_mode",
			Signature: "teach_mode()",
			Module:    "motion",
			Doc:       "Deprecated: Set robot in freedrive mode. Use freedrive_mode instead.",
		},
		{
			Name:      "track_conveyor_circular",
			Signature: "track_conveyor_circular(center, ticks_per_revolution, rotate_tool=False, encoder_index=0)",
			Module:    "motion",
			Doc: "Make robot movements track a circular conveyor.\n" +
				"@param center pose Center of the conveyor in base coordinates\n" +
				"@param ticks_per_revolution float Encoder ticks per full revolution\n" +
				"@param rotate_tool bool Whether the tool rotates with the conveyor\n" +
				"@param encoder_index int Encoder index, must be 0 or 1",
		},
		{
			Name:      "track_conveyor_linear",
			Signature: "track_conveyor_linear(direction, ticks_per_meter, encoder_index=0)",
			Module:    "motion",
			Doc: "Make robot movements track a linear conveyor.\n" +
				"@param direction pose Direction of the conveyor in base coordinates\n" +
				"@param ticks_per_meter float Encoder ticks per meter of conveyor travel\n" +
				"@param encoder_index int Encoder index, must be 0 or 1",
		},
		{
			Name:      "freedrive_mode_no_incorrect_payload_check",
			Signature: "freedrive_mode_no_incorrect_payload_check()",
			Module:    "motion",
			Doc:       "Enter freedrive mode without payload verification. User must ensure correct payload.",
		},
	}
}
