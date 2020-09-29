import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-variable-edit',
  templateUrl: './variable-edit.component.html',
  styleUrls: ['./variable-edit.component.scss']
})
export class VariableEditComponent implements OnInit {

  @Input() variable: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.variable.id) {
      this.as.get('variable/' + this.variable.id).subscribe(res => {
        this.variable = res.data;
      });
    }
  }

  submit(): void {
    if (this.variable.id) {
      this.as.put('variable/' + this.variable.id, this.variable).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('variable', this.variable).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
