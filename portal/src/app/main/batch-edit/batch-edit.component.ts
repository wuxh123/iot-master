import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-batch-edit',
  templateUrl: './batch-edit.component.html',
  styleUrls: ['./batch-edit.component.scss']
})
export class BatchEditComponent implements OnInit {

  @Input() batch: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.batch.id) {
      this.as.get('batch/' + this.batch.id).subscribe(res => {
        this.batch = res.data;
      });
    }
  }

  submit(): void {
    if (this.batch.id) {
      this.as.put('batch/' + this.batch.id, this.batch).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('batch', this.batch).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
