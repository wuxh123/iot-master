import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-strategy-edit',
  templateUrl: './strategy-edit.component.html',
  styleUrls: ['./strategy-edit.component.scss']
})
export class StrategyEditComponent implements OnInit {

  @Input() strategy: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.strategy.id) {
      this.as.get('strategy/' + this.strategy.id).subscribe(res => {
        this.strategy = res.data;
      });
    }
  }

  submit(): void {
    if (this.strategy.id) {
      this.as.put('strategy/' + this.strategy.id, this.strategy).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('strategy', this.strategy).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
