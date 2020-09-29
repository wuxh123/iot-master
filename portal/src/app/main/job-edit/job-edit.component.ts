import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-job-edit',
  templateUrl: './job-edit.component.html',
  styleUrls: ['./job-edit.component.scss']
})
export class JobEditComponent implements OnInit {

  @Input() job: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.job.id) {
      this.as.get('job/' + this.job.id).subscribe(res => {
        this.job = res.data;
      });
    }
  }

  submit(): void {
    if (this.job.id) {
      this.as.put('job/' + this.job.id, this.job).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('job', this.job).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
