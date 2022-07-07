module jarvis

go 1.18

replace task_manager => ./task_manager

replace media_manager => ./media_manager

replace monitor => ./monitor

require (
	media_manager v0.0.0-00010101000000-000000000000
	task_manager v0.0.0-00010101000000-000000000000
	monitor v0.0.0-00010101000000-000000000000
)
