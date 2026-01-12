package workerpool

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	p.metrics.IncWorkers()
	defer p.metrics.DecWorkers()

	for job := range p.jobs {
		select {
		case <-job.Ctx.Done():
			// job already abandoned
			continue
		default:
		}

		result := job.Do(job.Ctx)

		select {
		case job.ResultCh <- result:
			//delivered
		default:
			// handler gone drop result
		}

	}
}