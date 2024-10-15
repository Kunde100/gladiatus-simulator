<?php

require_once('simulate_arena.php');

header('Content-Type: application/json');

$request_method = $_SERVER['REQUEST_METHOD'];

if ($request_method == 'POST') {
	$input = json_decode(file_get_contents('php://input'), true);

	if (!isset($input['attacker']) || !isset($input['defender'])) {
		echo json_encode(array('error' => 'attacker and defender must be specified'));
		exit;
	}

	$options = isset($input['options']) ? $input['options'] : [];

	$results = arena_simulator($input['attacker'], $input['defender'], $options);

	echo json_encode($results);
	exit;
}

echo json_encode(['error' => true, 'message' => 'Invalid request method']);
?>
