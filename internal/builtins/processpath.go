package builtins

func processpathBuiltins() []Builtin {
	return []Builtin{
		{
			Name:      "mc_add_circular",
			Signature: "mc_add_circular(pose_via, pose_to, a, v, r, mode=0)",
			Module:    "processpath",
			Doc: "Queue a circular motion in the Part Coordinate System.\n" +
				"@param pose_via pose Via point on the arc (only position used) in PCS\n" +
				"@param pose_to pose Target pose in PCS\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param r float Blend radius [m]\n" +
				"@param mode int Orientation mode (0=unconstrained, 1=fixed relative to arc tangent)\n" +
				"@returns int Motion ID",
		},
		{
			Name:      "mc_add_linear",
			Signature: "mc_add_linear(pose, a, v, r)",
			Module:    "processpath",
			Doc: "Queue a linear motion in the Part Coordinate System.\n" +
				"@param pose pose Target pose in PCS\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s]\n" +
				"@param r float Blend radius [m]\n" +
				"@returns int Motion ID",
		},
		{
			Name:      "mc_add_path",
			Signature: "mc_add_path(path_id, a, v, r)",
			Module:    "processpath",
			Doc: "Queue motion from a loaded .nc toolpath file.\n" +
				"@param path_id int Path ID from mc_load_path\n" +
				"@param a float Tool acceleration [m/s^2]\n" +
				"@param v float Tool speed [m/s] (ignored if useFeedRate was True in mc_load_path)\n" +
				"@param r float Blend radius [m]\n" +
				"@returns int Motion ID",
		},
		{
			Name:      "mc_get_target_rtcp_speed",
			Signature: "mc_get_target_rtcp_speed()",
			Module:    "processpath",
			Doc: "Returns the target speed of a Remote TCP. For regular TCP use get_target_tcp_speed_along_path.\n" +
				"@returns float Speed in m/s (scalar)",
		},
		{
			Name:      "mc_initialize",
			Signature: "mc_initialize(mode, tcp, doc=6)",
			Module:    "processpath",
			Doc: "Initialize a new motion sequence. Resets speed factor, PCS, and motion ID counter.\n" +
				"@param mode int 0=regular TCP motion, 1=remote TCP motion\n" +
				"@param tcp pose TCP offset (regular mode) or remote TCP pose (remote mode)\n" +
				"@param doc int Degrees of constraint (6=all DOF, 5=free rotation about Z)",
		},
		{
			Name:      "mc_load_path",
			Signature: "mc_load_path(nc_file, useFeedRate)",
			Module:    "processpath",
			Doc: "Load a .nc toolpath file for use with mc_add_path.\n" +
				"@param nc_file string Path to the .nc file on disk\n" +
				"@param useFeedRate bool Whether to use feedrate from the NC file\n" +
				"@returns int Path ID for use with mc_add_path",
		},
		{
			Name:      "mc_run_motion",
			Signature: "mc_run_motion(id=-1)",
			Module:    "processpath",
			Doc: "Execute queued motions. Waits for the specified motion to complete.\n" +
				"@param id int Motion ID to wait for (-1 = run all queued motions)",
		},
		{
			Name:      "mc_set_pcs",
			Signature: "mc_set_pcs(pcs)",
			Module:    "processpath",
			Doc: "Set the Part Coordinate System for subsequent motions.\n" +
				"@param pcs pose Part Coordinate System relative to robot base",
		},
		{
			Name:      "mc_set_speed_factor",
			Signature: "mc_set_speed_factor(s)",
			Module:    "processpath",
			Doc: "Set a speed scaling factor for the current motion sequence.\n" +
				"@param s float Speed factor (1.0 = normal speed)",
		},
	}
}
